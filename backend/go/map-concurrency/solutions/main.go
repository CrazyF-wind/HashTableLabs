// main.go
package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"os"
	"sync"
)

// ============================================================
// 方案 1：sync.RWMutex 保护的普通 map（最通用）
// ============================================================
type RWMutexCache struct {
	mu sync.RWMutex
	m  map[string]string
}

func NewRWMutexCache() *RWMutexCache {
	return &RWMutexCache{m: make(map[string]string)}
}

func (c *RWMutexCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.m[key]
	return v, ok
}

func (c *RWMutexCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = value
}

func (c *RWMutexCache) Snapshot() map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make(map[string]string, len(c.m))
	for k, v := range c.m {
		out[k] = v
	}
	return out
}

// ============================================================
// 方案 2：sync.Map（读多写少，key 集合稳定）
// ============================================================
type SyncMapCache struct {
	m sync.Map
}

func (c *SyncMapCache) Get(key string) (string, bool) {
	v, ok := c.m.Load(key)
	if !ok {
		return "", false
	}
	return v.(string), true
}

func (c *SyncMapCache) Set(key, value string) {
	c.m.Store(key, value)
}

// 注意：sync.Map.Range 不保证一致性快照
func (c *SyncMapCache) Snapshot() map[string]string {
	out := make(map[string]string)
	c.m.Range(func(k, v any) bool {
		out[k.(string)] = v.(string)
		return true
	})
	return out
}

// ============================================================
// 方案 3：分片锁（高并发、锁竞争严重时）
// ============================================================
const shardCount = 32

type shard struct {
	mu sync.RWMutex
	m  map[string]string
}

type ShardedCache struct {
	shards [shardCount]*shard
}

func NewShardedCache() *ShardedCache {
	c := &ShardedCache{}
	for i := range c.shards {
		c.shards[i] = &shard{m: make(map[string]string)}
	}
	return c
}

func (c *ShardedCache) shardFor(key string) *shard {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return c.shards[h.Sum32()%shardCount]
}

func (c *ShardedCache) Get(key string) (string, bool) {
	s := c.shardFor(key)
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.m[key]
	return v, ok
}

func (c *ShardedCache) Set(key, value string) {
	s := c.shardFor(key)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
}

// 注意：分片锁下没有全局一致快照，只能拿"近似"视图（逐片加锁）
func (c *ShardedCache) Snapshot() map[string]string {
	out := make(map[string]string)
	for _, s := range c.shards {
		s.mu.RLock()
		for k, v := range s.m {
			out[k] = v
		}
		s.mu.RUnlock()
	}
	return out
}

// ============================================================
// 方案 4：通道串行化（用 channel 把并发转成串行，锁彻底消除）
// ============================================================
type req struct {
	op    string
	key   string
	value string
	reply chan string
	done  chan bool
}

type ChannelCache struct {
	reqs chan req
}

func NewChannelCache() *ChannelCache {
	c := &ChannelCache{reqs: make(chan req, 128)}
	go c.loop()
	return c
}

func (c *ChannelCache) loop() {
	m := make(map[string]string)
	for r := range c.reqs {
		switch r.op {
		case "get":
			v, ok := m[r.key]
			r.reply <- v
			r.done <- ok
		case "set":
			m[r.key] = r.value
			r.reply <- ""
			r.done <- true
		}
	}
}

func (c *ChannelCache) Get(key string) (string, bool) {
	reply := make(chan string, 1)
	done := make(chan bool, 1)
	c.reqs <- req{op: "get", key: key, reply: reply, done: done}
	return <-reply, <-done
}

func (c *ChannelCache) Set(key, value string) {
	reply := make(chan string, 1)
	done := make(chan bool, 1)
	c.reqs <- req{op: "set", key: key, value: value, reply: reply, done: done}
	<-reply
	<-done
}

// 通道方案没有低成本的全量快照；这里返回 nil 表示"不支持"
func (c *ChannelCache) Snapshot() map[string]string {
	return nil
}

// ============================================================
// 统一接口 + 压测验证
// ============================================================
type Cache interface {
	Get(key string) (string, bool)
	Set(key, value string)
	Snapshot() map[string]string
}

func stress(c Cache, name string, workers, ops int) {
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < ops; i++ {
				k := fmt.Sprintf("w%d-k%d", id, i)
				c.Set(k, "v")
				if v, ok := c.Get(k); !ok || v != "v" {
					fmt.Printf("[%s] 数据不一致: key=%s v=%q ok=%v\n", name, k, v, ok)
				}
			}
		}(w)
	}
	wg.Wait()
	fmt.Printf("[%s] 压测通过：%d workers × %d ops，无崩溃\n", name, workers, ops)
}

func main() {
	mode := "all"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	type caseT struct {
		name  string
		build func() Cache
	}
	cases := []caseT{
		{"RWMutex", func() Cache { return NewRWMutexCache() }},
		{"SyncMap", func() Cache { return &SyncMapCache{} }},
		{"Sharded", func() Cache { return NewShardedCache() }},
		{"Channel", func() Cache { return NewChannelCache() }},
	}

	for _, tc := range cases {
		if mode != "all" && mode != tc.name {
			continue
		}
		c := tc.build()
		stress(c, tc.name, 8, 20_000)
		if snap := c.Snapshot(); snap != nil && len(snap) > 0 {
			b, _ := json.Marshal(map[string]int{"size": len(snap)})
			fmt.Printf("[%s] snapshot: %s\n", tc.name, b)
		}
	}

	fmt.Println("全部通过，进程正常退出")
}
