# array - 线程安全动态数组库 / Thread-Safe Dynamic Array Utilities

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`array` 是一个基于 Go 语言的线程安全动态数组库，支持并发环境下的安全操作，包括添加、删除、插入、搜索等常见操作。适用于高并发场景下的数据集合管理。

GitHub地址: [github.com/small-ek/antgo/container/array](https://github.com/small-ek/antgo/container/array)

### 📦 安装

```bash
go get github.com/small-ek/antgo/container/array
```

### 🚀 快速开始

#### 初始化数组
```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/container/array"
)

func main() {
	// 创建线程安全数组（初始容量为10）
	arr := array.New[int](10)
}
```

#### 添加元素
```go
func main() {
	arr := array.New[int](0)
	arr.Append(1)       // 添加单个元素
	arr.Append(2, 3, 4) // 批量添加
	fmt.Println(arr.List()) // 输出 [1 2 3 4]
}
```

#### 删除元素
```go
func main() {
	arr := array.New[string](5)
	arr.Append("A", "B", "C", "D")

	// 删除索引为1的元素
	err := arr.Delete(1)
	if err != nil {
		fmt.Println("删除失败:", err)
	}
	fmt.Println(arr.List()) // 输出 [A C D]
}
```

#### 插入元素
```go
func main() {
	arr := array.New[float64](3)
	arr.Append(1.1, 3.3)

	// 在索引1处插入2.2
	err := arr.Insert(1, 2.2)
	if err != nil {
		fmt.Println("插入失败:", err)
	}
	fmt.Println(arr.List()) // 输出 [1.1 2.2 3.3]
}
```

### 🔧 高级用法

#### 并发安全遍历
```go
func main() {
	arr := array.New[int](10)
	arr.Append(1, 2, 3, 4, 5)

	// 读锁保护下的遍历
	arr.WithReadLock(func(data []int) {
		for _, v := range data {
			fmt.Println(v)
		}
	})
}
```

#### 批量操作
```go
func main() {
	arr := array.New[string](5)
	arr.Append("Apple", "Banana", "Cherry")

	// 写锁保护下的批量更新
	arr.WithWriteLock(func(data []string) {
		for i := range data {
			data[i] = data[i] + "_new"
		}
	})
	fmt.Println(arr.List()) // 输出 [Apple_new Banana_new Cherry_new]
}
```

### ✨ 核心特性

| 特性                | 描述                                                                 |
|---------------------|--------------------------------------------------------------------|
| **线程安全**         | 基于 `sync.RWMutex` 实现并发安全的读写操作                          |
| **泛型支持**         | 支持任意可比较类型（Go 1.18+）                                      |
| **高性能**           | 内存预分配与批量操作优化，减少锁竞争                                |
| **丰富的 API**       | 提供 `Append`、`Delete`、`Insert`、`Search` 等 10+ 种操作方法       |

### ⚠️ 注意事项
1. 使用 `Insert` 或 `Delete` 时需检查索引合法性，否则返回预定义错误 `ErrIndexOutOfBounds`。
2. 读操作（如 `Get`、`List`）使用读锁，写操作（如 `Append`、`Delete`）使用写锁。
3. 批量操作建议使用 `WithReadLock`/`WithWriteLock` 方法减少锁粒度。

### 🤝 参与贡献
[贡献指南](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [提交Issue](https://github.com/small-ek/antgo/issues)

---

## English

### 📖 Introduction

`array` is a Go-based thread-safe dynamic array library that supports common operations like appending, deleting, inserting, and searching in concurrent environments. Ideal for managing data collections in high-concurrency scenarios.

GitHub URL: [github.com/small-ek/antgo/container/array](https://github.com/small-ek/antgo/container/array)

### 📦 Installation

```bash
go get github.com/small-ek/antgo/container/array
```

### 🚀 Quick Start

#### Initialize Array
```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/container/array"
)

func main() {
	// Create a thread-safe array (initial capacity 10)
	arr := array.New[int](10)
}
```

#### Append Elements
```go
func main() {
	arr := array.New[int](0)
	arr.Append(1)       // Add single element
	arr.Append(2, 3, 4) // Batch append
	fmt.Println(arr.List()) // Output [1 2 3 4]
}
```

#### Delete Elements
```go
func main() {
	arr := array.New[string](5)
	arr.Append("A", "B", "C", "D")

	// Delete element at index 1
	err := arr.Delete(1)
	if err != nil {
		fmt.Println("Delete failed:", err)
	}
	fmt.Println(arr.List()) // Output [A C D]
}
```

#### Insert Elements
```go
func main() {
	arr := array.New[float64](3)
	arr.Append(1.1, 3.3)

	// Insert 2.2 at index 1
	err := arr.Insert(1, 2.2)
	if err != nil {
		fmt.Println("Insert failed:", err)
	}
	fmt.Println(arr.List()) // Output [1.1 2.2 3.3]
}
```

### 🔧 Advanced Usage

#### Concurrent-Safe Iteration
```go
func main() {
	arr := array.New[int](10)
	arr.Append(1, 2, 3, 4, 5)

	// Iterate under read lock
	arr.WithReadLock(func(data []int) {
		for _, v := range data {
			fmt.Println(v)
		}
	})
}
```

#### Batch Operations
```go
func main() {
	arr := array.New[string](5)
	arr.Append("Apple", "Banana", "Cherry")

	// Batch update under write lock
	arr.WithWriteLock(func(data []string) {
		for i := range data {
			data[i] = data[i] + "_new"
		}
	})
	fmt.Println(arr.List()) // Output [Apple_new Banana_new Cherry_new]
}
```

### ✨ Key Features

| Feature             | Description                                                     |
|---------------------|-----------------------------------------------------------------|
| **Thread-Safe**     | Implements concurrency-safe operations via `sync.RWMutex`      |
| **Generics**        | Supports any comparable type (Go 1.18+)                        |
| **High Performance**| Optimized with memory pre-allocation and batch operations      |
| **Rich API**        | Provides 10+ methods like `Append`, `Delete`, `Insert`, `Search` |

### ⚠️ Important Notes
1. Check index validity when using `Insert` or `Delete` to avoid `ErrIndexOutOfBounds`.
2. Read operations (e.g., `Get`, `List`) use read locks; write operations use write locks.
3. Use `WithReadLock`/`WithWriteLock` for batch operations to minimize lock contention.

### 🤝 Contributing
[Contribution Guide](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/small-ek/antgo/issues)

[⬆ Back to Top](#中文)