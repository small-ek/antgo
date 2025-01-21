```markdown
# antgo/encoding/abase64 - Base64 Encoding/Decoding Library

[中文](#中文) | [English](#English)

## 中文

### 简介

`antgo/encoding/abase64` 是一个高效的Base64编码解码库，基于Go语言标准库实现，优化了内存分配和性能。  
支持标准Base64编码，适用于处理敏感数据、文件编码或网络传输场景。

GitHub地址: [github.com/small-ek/antgo/encoding/abase64](https://github.com/small-ek/antgo/encoding/abase64)

### 安装

使用Go Modules安装：

```bash
go get github.com/small-ek/antgo/encoding/abase64
```

### 使用示例

#### 编码

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/encoding/abase64"
)

func main() {
	data := []byte("Hello, World!")
	encoded := abase64.Encode(data)
	fmt.Println(encoded) // 输出: SGVsbG8sIFdvcmxkIQ==
}
```

#### 解码

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/encoding/abase64"
)

func main() {
	encodedStr := "SGVsbG8sIFdvcmxkIQ=="
	decoded, err := abase64.Decode(encodedStr)
	if err != nil {
		fmt.Println("解码错误:", err)
		return
	}
	fmt.Println(string(decoded)) // 输出: Hello, World!
}
```

### 特点

- **高效内存管理**: 预计算缓冲区大小，减少内存分配次数。
- **符合RFC标准**: 使用`base64.StdEncoding`，兼容性高。
- **错误处理**: 解码时自动验证输入合法性。

### 注意事项

- 输入字符串必须为标准Base64格式，否则解码会返回错误。
- 支持`+`、`/`字符，若需URL安全版本可后续扩展。

### 贡献

欢迎提交Issue或PR: [贡献指南](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md)

---

## English

### Introduction

`antgo/encoding/abase64` is a high-performance Base64 encoding/decoding library based on Go's standard library, optimized for memory allocation and performance.  
It supports standard Base64 encoding, suitable for sensitive data handling, file encoding, or network transmission.

GitHub URL: [github.com/small-ek/antgo/encoding/abase64](https://github.com/small-ek/antgo/encoding/abase64)

### Installation

Using Go Modules:

```bash
go get github.com/small-ek/antgo/encoding/abase64
```

### Usage Examples

#### Encoding

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/encoding/abase64"
)

func main() {
	data := []byte("Hello, World!")
	encoded := abase64.Encode(data)
	fmt.Println(encoded) // Output: SGVsbG8sIFdvcmxkIQ==
}
```

#### Decoding

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/encoding/abase64"
)

func main() {
	encodedStr := "SGVsbG8sIFdvcmxkIQ=="
	decoded, err := abase64.Decode(encodedStr)
	if err != nil {
		fmt.Println("Decode error:", err)
		return
	}
	fmt.Println(string(decoded)) // Output: Hello, World!
}
```

### Features

- **Efficient Memory**: Pre-calculates buffer size to minimize allocations.
- **RFC-Compliant**: Uses `base64.StdEncoding` for compatibility.
- **Error Handling**: Validates input during decoding.

### Notes

- Input must be standard Base64; invalid formats return errors.
- Supports `+` and `/` characters. Contact us if URL-safe version is needed.

### Contributing

Issues and PRs are welcome: [Contribution Guide](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md)
