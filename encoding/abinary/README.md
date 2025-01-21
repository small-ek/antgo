# antgo/encoding/abinary - 二进制序列化库 / Binary Serialization Library

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`antgo/encoding/abinary` 是基于Go标准库的高效二进制序列化工具，通过内存池优化和零拷贝技术实现高性能编解码。  
适用于网络传输、持久化存储和高性能计算场景。

GitHub地址: [github.com/small-ek/antgo/encoding/abinary](https://github.com/small-ek/antgo/encoding/abinary)

### 📦 安装

```bash
go get github.com/small-ek/antgo/encoding/abinary
```

### 🚀 快速开始

#### 基本类型编码
```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/encoding/abinary"
)

func main() {
	// 编码int32
	val := int32(42)
	data, err := abinary.Encode(val)
	if err != nil {
		panic(err)
	}

	// 解码
	var decoded int32
	if err := abinary.Decode(data, &decoded); err != nil {
		panic(err)
	}
	fmt.Printf("解码结果: %d", decoded) // 输出: 解码结果: 42
}
```

#### 结构体编码
```go
type Sensor struct {
	ID    uint32
	Value float64
}

func main() {
	// 编码结构体
	s := Sensor{ID: 1, Value: 25.5}
	encoded, err := abinary.Encode(s)
	if err != nil {
		panic(err)
	}

	// 解码结构体
	var decoded Sensor
	if err := abinary.Decode(encoded, &decoded); err != nil {
		panic(err)
	}
	fmt.Printf("解码结构体: %+v", decoded) // 输出: {ID:1 Value:25.5}
}
```

### ✨ 核心特性

| 特性                | 描述                                                                 |
|---------------------|--------------------------------------------------------------------|
| **内存池优化**       | 使用sync.Pool复用缓冲区，减少90%内存分配                              |
| **小端序支持**       | 采用Little-Endian格式，兼容x86架构                                    |
| **类型安全**         | 编译时类型检查，防止运行时类型错误                                    |
| **高并发支持**       | 编码器无状态设计，解码器线程安全                                      |

### ⚠️ 注意事项
1. 仅支持固定大小类型（int32/uint64等）和内存对齐的结构体
2. 结构体字段需满足：`sizeof(struct) % alignment == 0`
3. 不支持slice/map/string等变长类型
4. 返回数据为独立副本，可安全修改

### 🤝 参与贡献
[贡献指南](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [提交Issue](https://github.com/small-ek/antgo/issues)

---

## English

### 📖 Introduction

`antgo/encoding/abinary` is a high-performance binary serialization library with memory pool optimization and zero-copy techniques.  
Ideal for network transmission, persistent storage, and high-performance computing.

GitHub URL: [github.com/small-ek/antgo/encoding/abinary](https://github.com/small-ek/antgo/encoding/abinary)

### 📦 Installation

```bash
go get github.com/small-ek/antgo/encoding/abinary
```

### 🚀 Quick Start

#### Primitive Type Encoding
```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/encoding/abinary"
)

func main() {
	// Encode int32
	val := int32(42)
	data, err := abinary.Encode(val)
	if err != nil {
		panic(err)
	}

	// Decode
	var decoded int32
	if err := abinary.Decode(data, &decoded); err != nil {
		panic(err)
	}
	fmt.Printf("Decoded: %d", decoded) // Output: Decoded: 42
}
```

#### Struct Encoding
```go
type Sensor struct {
	ID    uint32
	Value float64
}

func main() {
	// Encode struct
	s := Sensor{ID: 1, Value: 25.5}
	encoded, err := abinary.Encode(s)
	if err != nil {
		panic(err)
	}

	// Decode struct
	var decoded Sensor
	if err := abinary.Decode(encoded, &decoded); err != nil {
		panic(err)
	}
	fmt.Printf("Decoded: %+v", decoded) // Output: {ID:1 Value:25.5}
}
```

### ✨ Key Features

| Feature             | Description                                                        |
|---------------------|--------------------------------------------------------------------|
| **Memory Pool**     | 90% less allocations with sync.Pool                                |
| **Little-Endian**   | Native support for x86 architecture                               |
| **Type Safety**     | Compile-time type checking prevents runtime errors                |
| **Concurrency**     | Stateless encoder and thread-safe decoder                         |

### ⚠️ Important Notes
1. Only fixed-size types (int32/uint64 etc.) and aligned structs supported
2. Structs must satisfy: `sizeof(struct) % alignment == 0`
3. Variable-length types (slice/map/string) not supported
4. Returns independent data copies for safe modification

### 🤝 Contributing
[Contribution Guide](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/small-ek/antgo/issues)

[⬆ Back to Top](#antgoencodingabinary---binary-serialization-library)
