English | [中文](README_ZH.md)

# 📚 Introduction

`ring_buf` It is a high-performance ring cache built from an array.

`ring_buf` It integrates functions such as array expansion, shrinkage, reading without data hibernation, etc.

# 🚀 Features

- [x] [High-performance](#-performance) Based on the principle of write priority to ensure high throughput
- [x] built-in memory pool
- [x] Integrate `ring_buf` to expand and shrink, efficient, reusable and auto-scalable memory
- [x] State control read and write, make read and write execution ☞ a line

# 🎬 Getting started

`ring_buf` is available as a Go module and we highly recommend that you use `ring_buf` via [Go Modules](https://go.dev/blog/using-go-modules), with Go 1.11 Modules enabled (Go 1.11+), you can just simply add `import "github.com/byfengfeng/ring_buf"` to the codebase and run `go mod download/go mod tidy` or `go [build|run|test]` to download the necessary dependencies automatically.

## With v1

```powershell
go get -u github.com/byfengfeng/ring_buf
```
