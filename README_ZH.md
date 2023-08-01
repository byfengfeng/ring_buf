[英文](README.md) | 中文

# 📚 简介

`ring_buf` 是一个数组进行构建的高性能环形缓存。

`ring_buf` 集成了数组扩容、缩容、读取无数据休眠等功能。

# 🚀 功能

- [x] [高性能](#-性能测试) 的基于写入优先原则，保证高吞吐
- [x] 内置内存池
- [x] 集成`ring_buf`扩容、缩容，高效、可重用而且自动伸缩的内存
- [x] 状态控制读写，使读写执行☞一条线上
# 🎬 开始

`ring_buf` 是一个 Go module，而且我们也强烈推荐通过 [Go Modules](https://go.dev/blog/using-go-modules) 来使用 `ring_buf`，在开启 Go Modules 支持（Go 1.19+）之后可以通过简单地在代码中写 `import "github.com/byfengfeng/ring_buf"` 来引入 `gnet`，然后执行 `go mod download/go mod tidy` 或者 `go [build|run|test]` 这些命令来自动下载所依赖的包。

## 使用 v1

```powershell
go get -u github.com/byfengfeng/ring_buf
```
