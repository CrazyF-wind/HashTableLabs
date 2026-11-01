# 一次 map 并发写崩溃

对应文章：[一次 map 并发写崩溃，让我重新理解了 Go runtime](https://crazyf-wind.github.io/2026/10/26/一次%20map%20并发写崩溃，让我重新理解了%20Go%20runtime/)

## 目录

- `reproduce/`：复现 Go 原生 map 在并发读写时的崩溃
- `solutions/`：四种并发安全方案的压测对比

## 复现崩溃

```bash
cd backend/go/map-concurrency

# 单 goroutine 顺序访问，跑 100 万次也不会崩
go run ./reproduce sequential

# 8 个 goroutine 并发写，通常会以 fatal error 退出
go run ./reproduce concurrent

# 先顺序、再并发，对比同一份业务代码
go run ./reproduce compare
```

崩溃时日志里的 `fatal error: concurrent map writes` 来自 runtime 的 `throw`，不能被 `recover`。用 `go run` 运行时外层工具返回 1，程序本身的退出码是 2。

## 四种方案

```bash
go run ./solutions            # 全部
go run ./solutions RWMutex    # 只跑某一种
go run ./solutions SyncMap
go run ./solutions Sharded
go run ./solutions Channel
```

## 用 race detector 定位

```bash
go run -race ./reproduce concurrent
```
