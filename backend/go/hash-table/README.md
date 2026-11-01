# 从零实现一个哈希表

对应文章：[从零实现一个哈希表，理解它的每一个设计决策](https://crazyf-wind.github.io/2026/10/30/从零实现一个哈希表，理解它的每一个设计决策/)

链地址法实现：FNV-1a 哈希、桶数为 2 的幂、负载因子超过 0.75 翻倍扩容、`Range` 遍历、不提供并发安全。

## 文件

- `hashmap.go`：完整实现
- `hashmap_test.go`：基础用例
- `hashmap_extra_test.go`：白盒、边界与差分测试

## 运行

```bash
cd backend/go/hash-table
go test ./...
go test -race ./...
```

指定用例：

```bash
go test -run TestDifferentialAgainstNativeMap -v ./...
```

注意：不要用 `go test hashmap_test.go`。显式列出文件时 Go 只编译点名的文件，`hashmap.go` 不参与，会报 `New` 未定义；把目录交给 `./...` 或 `.` 即可。
