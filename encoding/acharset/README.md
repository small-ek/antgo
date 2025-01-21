# antgo/encoding/acharset - 字符集编解码库 / Charset Encoding Library

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`antgo/encoding/acharset` 是基于Go标准库的高效字符集编解码工具，支持多种字符集别名映射，并通过并发安全缓存优化编解码性能。  
适用于处理多国语言文本编码转换、旧系统数据迁移等场景。

GitHub地址: [github.com/small-ek/antgo/encoding/acharset](https://github.com/small-ek/antgo/encoding/acharset)

### 📦 安装

```bash
go get github.com/small-ek/antgo/encoding/acharset
```

### 🚀 快速开始

#### 基本编解码
```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/encoding/acharset"
)

func main() {
	// GBK编码的原始字节
	gbkBytes := []byte{0xC4, 0xE3, 0xBA, 0xC3} // "你好"的GBK编码
	
	// 解码为UTF-8
	utf8Bytes, err := acharset.Decode(string(gbkBytes), "GBK")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(utf8Bytes)) // 输出: 你好
}
```

#### 处理HZ-GB-2312编码
```go
func main() {
	// HZ-GB-2312编码文本
	hzText := "~{<:Ky2;S{<~}" // "你好世界"的HZ编码
	
	// 转换为UTF-8
	utf8Bytes, err := acharset.Decode(hzText, "hzgb2312")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(utf8Bytes)) // 输出: 你好世界
}
```

### ✨ 核心特性

| 特性                | 描述                                                                 |
|---------------------|--------------------------------------------------------------------|
| **别名支持**         | 内置常见字符集别名（如GB2312→HZ-GB-2312）                           |
| **缓存优化**         | 使用sync.Map缓存已解析的编码器，提升重复使用性能                     |
| **并发安全**         | 所有操作线程安全，适合高并发场景                                     |
| **自动规范化**       | 自动处理字符集名称大小写（如"gbk"→"GBK"）                           |

### ⚠️ 注意事项
1. 支持的字符集取决于系统环境中的IANA注册表
2. 内置别名可通过修改`charsetAliases`扩展
3. 解码失败会返回`unsupported charset`错误
4. 返回的字节切片是独立副本，可安全修改

### 🤝 参与贡献
[贡献指南](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [提交Issue](https://github.com/small-ek/antgo/issues)

---

## English

### 📖 Introduction

`antgo/encoding/acharset` is an efficient character set encoding/decoding library with alias support and concurrent-safe caching.  
Ideal for multi-language text processing and legacy system data migration.

GitHub URL: [github.com/small-ek/antgo/encoding/acharset](https://github.com/small-ek/antgo/encoding/acharset)

### 📦 Installation

```bash
go get github.com/small-ek/antgo/encoding/acharset
```

### 🚀 Quick Start

#### Basic Encoding/Decoding
```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/encoding/acharset"
)

func main() {
	// GBK encoded bytes
	gbkBytes := []byte{0xC4, 0xE3, 0xBA, 0xC3} // "Hello" in GBK
	
	// Decode to UTF-8
	utf8Bytes, err := acharset.Decode(string(gbkBytes), "GBK")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(utf8Bytes)) // Output: 你好
}
```

#### Handling HZ-GB-2312
```go
func main() {
	// HZ-GB-2312 encoded text
	hzText := "~{<:Ky2;S{<~}" // "Hello world" in HZ encoding
	
	// Convert to UTF-8
	utf8Bytes, err := acharset.Decode(hzText, "hzgb2312")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(utf8Bytes)) // Output: 你好世界
}
```

### ✨ Key Features

| Feature             | Description                                                        |
|---------------------|--------------------------------------------------------------------|
| **Alias Support**   | Built-in charset aliases (e.g. GB2312→HZ-GB-2312)                 |
| **Caching**         | sync.Map cached encodings for repeated use                        |
| **Concurrency**     | Thread-safe operations for high concurrency scenarios             |
| **Auto-Normalize**  | Case-insensitive charset name handling (e.g. "gbk"→"GBK")         |

### ⚠️ Important Notes
1. Supported charsets depend on system's IANA registry
2. Extend aliases via modifying `charsetAliases`
3. Returns `unsupported charset` on decoding failure
4. Returned byte slice is an independent copy

### 🤝 Contributing
[Contribution Guide](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/small-ek/antgo/issues)

[⬆ Back to Top](#antgoencodingacharset---charset-encoding-library)

---