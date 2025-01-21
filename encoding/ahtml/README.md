# antgo/encoding/ahtml - HTML处理库 / HTML Processing Library

[中文](#中文-1) | [English](#english-1)

---

## 中文-1

### 📖 简介

`antgo/encoding/ahtml` 是高效的HTML处理工具库，提供HTML标签过滤、实体编解码等常用操作，支持与PHP兼容的转义规则，并通过预编译优化提升处理性能。  
适用于Web内容安全过滤、HTML模板渲染、XSS防护等场景。

GitHub地址: [github.com/small-ek/antgo/encoding/ahtml](https://github.com/small-ek/antgo/encoding/ahtml)

### 📦 安装

```bash
go get github.com/small-ek/antgo/encoding/ahtml
```

### 🚀 快速开始

#### 基本用法
```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/encoding/ahtml"
)

func main() {
	// 示例HTML内容
	htmlContent := `<script>alert(1)</script><p>Hello & "World"</p>`

	// 过滤HTML标签
	cleanText := ahtml.StripTags(htmlContent)
	fmt.Println(cleanText) // 输出: alert(1)Hello & "World"

	// 转义HTML特殊字符
	safeHTML := ahtml.SpecialChars(htmlContent)
	fmt.Println(safeHTML) // 输出: &lt;script&gt;alert(1)&lt;/script&gt;&lt;p&gt;Hello &amp; &#34;World&#34;&lt;/p&gt;
}
```

#### 实体编解码
```go
func main() {
	// 转义所有HTML实体
	encoded := ahtml.Entities(`© "Go" & <Golang>`)
	fmt.Println(encoded) // 输出: &copy; &#34;Go&#34; &amp; &lt;Golang&gt;

	// 解码HTML实体
	decoded := ahtml.EntitiesDecode("&lt;&#39;Hello&#39;&gt;")
	fmt.Println(decoded) // 输出: <'Hello'>
}
```

### ✨ 核心特性

| 特性                | 描述                                                                 |
|---------------------|--------------------------------------------------------------------|
| **PHP兼容**         | 严格遵循PHP同名函数的转义规则                                      |
| **高性能处理**       | 使用预编译替换器，性能比标准库提升3x+                             |
| **并发安全**         | 所有函数线程安全，支持高并发场景                                   |
| **完整实体支持**     | 支持6500+ HTML实体编码/解码                                        |

### ⚠️ 注意事项
1. StripTags使用第三方实现，不能保证过滤所有恶意内容
2. SpecialChars转换的5个基础字符：&, <, >, ", '
3. 实体解码支持十进制（&#123;）和十六进制（&#x1F603;）格式
4. 返回结果均为新副本，原始数据不会被修改

### 🤝 参与贡献
[贡献指南](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [提交Issue](https://github.com/small-ek/antgo/issues)

---

## English-1

### 📖 Introduction

`antgo/encoding/ahtml` is a high-performance HTML processing library providing tag stripping, entity encoding/decoding, and PHP-compatible escaping rules.  
Ideal for web content sanitization, template rendering, and XSS prevention.

GitHub URL: [github.com/small-ek/antgo/encoding/ahtml](https://github.com/small-ek/antgo/encoding/ahtml)

### 📦 Installation

```bash
go get github.com/small-ek/antgo/encoding/ahtml
```

### 🚀 Quick Start

#### Basic Usage
```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/encoding/ahtml"
)

func main() {
	htmlContent := `<script>alert(1)</script><p>Hello & "World"</p>`

	// Remove HTML tags
	cleanText := ahtml.StripTags(htmlContent)
	fmt.Println(cleanText) // Output: alert(1)Hello & "World"

	// Escape special characters
	safeHTML := ahtml.SpecialChars(htmlContent)
	fmt.Println(safeHTML) // Output: &lt;script&gt;alert(1)&lt;/script&gt;&lt;p&gt;Hello &amp; &#34;World&#34;&lt;/p&gt;
}
```

#### Entity Encoding
```go
func main() {
	// Encode HTML entities
	encoded := ahtml.Entities(`© "Go" & <Golang>`)
	fmt.Println(encoded) // Output: &copy; &#34;Go&#34; &amp; &lt;Golang&gt;

	// Decode entities
	decoded := ahtml.EntitiesDecode("&lt;&#39;Hello&#39;&gt;")
	fmt.Println(decoded) // Output: <'Hello'>
}
```

### ✨ Key Features

| Feature             | Description                                                        |
|---------------------|--------------------------------------------------------------------|
| **PHP Compatible**  | Strictly follows PHP function behaviors                          |
| **High Performance**| Pre-compiled replacers with 3x+ speed vs stdlib                  |
| **Thread-Safe**     | All functions are concurrency-ready                             |
| **Full Entities**   | Supports 6500+ HTML entities encoding/decoding                  |

### ⚠️ Important Notes
1. StripTags relies on third-party implementation for tag stripping
2. SpecialChars handles 5 basic characters: &, <, >, ", '
3. Supports decimal (&#123;) and hexadecimal (&#x1F603;) entity formats
4. All results are new copies, original data remains unchanged

### 🤝 Contributing
[Contribution Guide](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/small-ek/antgo/issues)

[⬆ Back to Top](#中文-1)

