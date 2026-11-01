// main.go
package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// ========== 业务代码：同一份，不改动 ==========
var cache = map[string]string{}

func Get(key string) string {
	return cache[key]
}

func Set(key, value string) {
	cache[key] = value
}

// ========== 场景 A：顺序单 goroutine 访问 ==========
func runSequential(n int) {
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("k%d", i)
		Set(key, "v")
		_ = Get(key)
	}
	fmt.Printf("[A] 顺序执行 %d 次读写，正常结束，cache size=%d\n", n, len(cache))
}

// ========== 场景 B：多 goroutine 并发访问 ==========
func runConcurrent(workers, ops int) {
	var wg sync.WaitGroup
	var counter int64

	start := make(chan struct{}) // 让所有 goroutine 尽量同时开始，提高撞上概率

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			<-start
			for i := 0; i < ops; i++ {
				key := fmt.Sprintf("w%d-k%d", id, i)
				Set(key, "v")
				_ = Get(key)
				atomic.AddInt64(&counter, 1)
			}
		}(w)
	}

	close(start)
	wg.Wait()
	fmt.Printf("[B] 并发执行完毕，ops=%d，正常结束（本次没撞上）\n", counter)
}

func main() {
	mode := "sequential"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}
	fmt.Printf("mode=%s GOMAXPROCS=%d\n", mode, runtime.GOMAXPROCS(0))

	switch mode {
	case "sequential":
		// 单 goroutine，跑 100 万次也不会崩
		runSequential(1_000_000)

	case "concurrent":
		// 多 goroutine 并发，几乎必崩
		runConcurrent(8, 100_000)

	case "compare":
		// 一次跑完对比：先顺序（安全），再并发（崩）
		fmt.Println("=== 第一阶段：单 goroutine 顺序访问 ===")
		runSequential(200_000)

		fmt.Println("=== 第二阶段：切换为多 goroutine 并发访问（同一份业务代码） ===")
		time.Sleep(50 * time.Millisecond)
		runConcurrent(8, 100_000)
		fmt.Println("如果这行被打印，说明本次侥幸没撞上，多跑几次即可复现")

	default:
		fmt.Println("usage: go run main.go [sequential|concurrent|compare]")
	}
}
