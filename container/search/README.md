# search - 并发安全搜索工具库 / Thread-Safe Search Utilities

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`search` 是一个高性能的Go语言泛型搜索工具库，提供线性搜索和二分搜索实现。支持所有可比较数据类型，专为并发场景设计，适合在多种数据集合中快速定位元素。

GitHub地址: [github.com/small-ek/antgo/container/search](https://github.com/small-ek/antgo/container/search)

### 📦 安装

```bash
go get github.com/small-ek/antgo/container/search
```

### 🚀 快速开始

#### 线性搜索
```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/container/search"
)

func main() {
	// 字符串切片搜索
	names := []string{"Alice", "Bob", "Charlie"}
	index := search.Search(names, "Bob")
	fmt.Println(index) // 输出: 1

	// 整型切片搜索
	numbers := []int{10, 20, 30, 15}
	fmt.Println(search.Search(numbers, 15)) // 输出: 3
}
```

#### 二分搜索
```go
func main() {
	// 必须有序的切片
	sorted := []int{1, 3, 5, 7, 9, 11}

	// 查找存在的元素
	index := search.SearchOrdered(sorted, 7)
	fmt.Println(index) // 输出: 3

	// 查找不存在的元素
	fmt.Println(search.SearchOrdered(sorted, 8)) // 输出: -1
}
```

### 🔧 高级用法

#### 自定义类型搜索
```go
type Product struct {
	ID    int
	Name  string
	Price float64
}

func main() {
	products := []Product{
		{101, "Keyboard", 29.99},
		{102, "Mouse", 19.95},
	}

	// 自定义相等判断
	index := search.Search(products, Product{ID: 102})
	fmt.Println(index) // 输出: 1
}
```

#### 并发环境使用
```go
func concurrentSearch() {
	data := []float64{1.1, 2.2, 3.3, 4.4, 5.5}

	// 多个goroutine并发搜索
	go func() {
		fmt.Println(search.Search(data, 3.3)) // 输出: 2
	}()

	go func() {
		fmt.Println(search.SearchOrdered(data, 4.4)) // 输出: 3
	}()
}
```

### ✨ 核心特性

| 特性                | 描述                                                                 |
|---------------------|--------------------------------------------------------------------|
| **双算法支持**       | 提供线性搜索(通用)和二分搜索(有序数据)两种模式                      |
| **泛型实现**         | 支持所有可比较类型(Go 1.18+)                                       |
| **零依赖**          | 仅使用标准库实现                                                   |
| **高性能**          | 二分搜索使用位运算优化                                             |
| **安全并发**        | 无状态设计，原生支持并发调用                                       |

### ⚠️ 注意事项
1. 使用`SearchOrdered`前必须确保切片**已按升序排列**
2. 自定义结构体类型需要实现`==`运算符
3. 二分搜索时间复杂度为O(log n)，线性搜索为O(n)
4. 未找到元素时统一返回-1

### 🤝 参与贡献
[贡献指南](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [提交Issue](https://github.com/small-ek/antgo/issues)

---

## English

### 📖 Introduction

`search` is a high-performance Go search utility library providing both linear and binary search implementations. Supports all comparable data types, designed for concurrent environments and fast element lookup in various collections.

GitHub URL: [github.com/small-ek/antgo/container/search](https://github.com/small-ek/antgo/container/search)

### 📦 Installation

```bash
go get github.com/small-ek/antgo/container/search
```

### 🚀 Quick Start

#### Linear Search
```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/container/search"
)

func main() {
	// String slice search
	names := []string{"Alice", "Bob", "Charlie"}
	index := search.Search(names, "Bob")
	fmt.Println(index) // Output: 1

	// Integer slice search
	numbers := []int{10, 20, 30, 15}
	fmt.Println(search.Search(numbers, 15)) // Output: 3
}
```

#### Binary Search
```go
func main() {
	// Must be sorted slice
	sorted := []int{1, 3, 5, 7, 9, 11}

	// Search existing element
	index := search.SearchOrdered(sorted, 7)
	fmt.Println(index) // Output: 3

	// Search non-existing element
	fmt.Println(search.SearchOrdered(sorted, 8)) // Output: -1
}
```

### 🔧 Advanced Usage

#### Custom Type Search
```go
type Product struct {
	ID    int
	Name  string
	Price float64
}

func main() {
	products := []Product{
		{101, "Keyboard", 29.99},
		{102, "Mouse", 19.95},
	}

	// Custom equality check
	index := search.Search(products, Product{ID: 102})
	fmt.Println(index) // Output: 1
}
```

#### Concurrent Usage
```go
func concurrentSearch() {
	data := []float64{1.1, 2.2, 3.3, 4.4, 5.5}

	// Concurrent searches in goroutines
	go func() {
		fmt.Println(search.Search(data, 3.3)) // Output: 2
	}()

	go func() {
		fmt.Println(search.SearchOrdered(data, 4.4)) // Output: 3
	}()
}
```

### ✨ Key Features

| Feature             | Description                                                     |
|---------------------|-----------------------------------------------------------------|
| **Dual Algorithms** | Provides both linear(universal) and binary(sorted data) search |
| **Generics**        | Supports all comparable types (Go 1.18+)                       |
| **Zero Dependency** | Implemented using standard library only                        |
| **High Performance**| Binary search optimized with bitwise operations                |
| **Concurrency Safe**| Stateless design for native concurrent calls                   |

### ⚠️ Important Notes
1. Slice **must be sorted in ascending order** before using `SearchOrdered`
2. Custom structs must implement `==` operator
3. Time complexity: O(log n) for binary, O(n) for linear search
4. Returns -1 uniformly when element not found

### 🤝 Contributing
[Contribution Guide](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/small-ek/antgo/issues)

[⬆ Back to Top](#中文)