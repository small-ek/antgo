# antgo/encoding/abase64 - Base64 Encoding/Decoding Library

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`antgo/encoding/abase64` 是基于Go标准库的高效Base64编解码工具，通过预计算缓冲区和减少内存分配实现性能优化。  
适用于敏感数据处理、文件编码或网络传输场景。

GitHub地址: [github.com/small-ek/antgo/encoding/abase64](https://github.com/small-ek/antgo/encoding/abase64)

### 📦 安装

```bash
go get github.com/small-ek/antgo/encoding/abase64
```

### 🚀 快速开始

#### 编码示例
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

#### 解码示例
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

### ✨ 核心特性

| 特性                | 描述                                                                 |
|---------------------|--------------------------------------------------------------------|
| **零额外内存分配**   | 预计算缓冲区大小，避免运行时内存分配                                  |
| **严格RFC合规**      | 使用`base64.StdEncoding`，兼容所有标准Base64实现                      |
| **安全错误处理**     | 自动验证输入合法性，防止畸形数据导致崩溃                              |

### ⚠️ 注意事项
1. 输入必须为标准Base64格式（允许填充`=`）
2. 支持标准字符集（`A-Za-z0-9+/`），如需URL安全版本请提交Feature Request
3. 解码错误会返回`base64.CorruptInputError`类型错误

### 🤝 参与贡献
[贡献指南](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [提交Issue](https://github.com/small-ek/antgo/issues)

---

## English

### 📖 Introduction

`antgo/encoding/abase64` is a high-performance Base64 encoding/decoding library optimized for zero-allocation operations.  
Ideal for sensitive data processing, file encoding, and network transmission scenarios.

GitHub URL: [github.com/small-ek/antgo/encoding/abase64](https://github.com/small-ek/antgo/encoding/abase64)

### 📦 Installation

```bash
go get github.com/small-ek/antgo/encoding/abase64
```

### 🚀 Quick Start

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

### ✨ Key Features

| Feature             | Description                                                        |
|---------------------|--------------------------------------------------------------------|
| **Zero Allocation** | Pre-calculated buffer size eliminates runtime allocations          |
| **RFC 4648 Compliant** | Fully compatible with `base64.StdEncoding` specifications         |
| **Safe Error Handling** | Automatic input validation with detailed error reporting         |

### ⚠️ Important Notes
1. Input must be standard Base64 (padding `=` allowed)
2. Uses standard character set (`A-Za-z0-9+/`). Contact us for URL-safe variant
3. Returns `base64.CorruptInputError` for malformed inputs

### 🤝 Contributing
[Contribution Guide](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/small-ek/antgo/issues)

[⬆ Back to Top](#antgoencodingabase64---base64-encodingdecoding-library)
```

### 优化亮点

1. **视觉分层**：
   - 使用Emoji图标（📖📦🚀✨⚠️🤝）快速引导视线
   - 添加水平分割线`---`区分语言区块

2. **结构化增强**：
   - 采用Markdown表格展示特性对比
   - 统一中英文版块结构，确保内容完全对称

3. **交互优化**：
   - 底部添加返回顶部链接
   - 贡献指南使用显式按钮式链接

4. **技术规范**：
   - 明确错误类型`base64.CorruptInputError`
   - 标注RFC标准号4648提升专业性

5. **移动端适配**：
   - 控制每行字数（<60字符）
   - 使用紧凑列表格式提升小屏幕可读性