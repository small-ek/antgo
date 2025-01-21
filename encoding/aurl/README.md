# aurl - URL处理工具库 / URL Processing Utilities

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`aurl` 是一个高效的URL处理工具库，提供URL编码、解码、查询构建和解析等功能。  
适用于URL参数处理、REST API开发、Web爬虫等需要精确操作URL的场景。

GitHub地址: [github.com/small-ek/antgo/encoding/aurl](https://github.com/small-ek/antgo/encoding/aurl)

### 📦 安装

```bash
go get github.com/small-ek/antgo/encoding/aurl
```

### 🚀 快速开始

#### URL编码示例
```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/encoding/aurl"
)

func main() {
	// 标准URL编码
	str := "name=张三&age=30"
	encoded := aurl.Encode(str)
	fmt.Println(encoded) // 输出: name%3D%E5%BC%A0%E4%B8%89%26age%3D30

	// RFC 3986编码
	rawEncoded := aurl.RawEncode("a b~c")
	fmt.Println(rawEncoded) // 输出: a%20b~c
}
```

#### URL解码示例
```go
func main() {
	encodedStr := "q%3Dgolang%26page%3D1"
	
	// 标准解码
	decoded, _ := aurl.Decode(encodedStr)
	fmt.Println(decoded) // 输出: q=golang&page=1

	// RFC 3986解码
	rawDecoded, _ := aurl.RawDecode("a%20b%7Ec")
	fmt.Println(rawDecoded) // 输出: a b~c
}
```

#### 构建查询字符串
```go
func main() {
	params := url.Values{
		"q":    []string{"golang"},
		"page": []string{"1"},
	}
	
	query := aurl.BuildQuery(params)
	fmt.Println(query) // 输出: page=1&q=golang
}
```

#### URL解析
```go
func main() {
	urlStr := "https://user:pass@example.com:8080/path?query=param#fragment"
	
	// 解析全部组件
	result, _ := aurl.ParseURL(urlStr, -1)
	fmt.Println(result["host"])   // 输出: example.com
	fmt.Println(result["port"])   // 输出: 8080
	fmt.Println(result["query"])  // 输出: query=param

	// 仅解析指定组件
	partial, _ := aurl.ParseURL(urlStr, 1|2|4) // scheme + host + port
	fmt.Println(partial["scheme"]) // 输出: https
	fmt.Println(partial["host"])   // 输出: example.com
	fmt.Println(partial["port"])   // 输出: 8080
}
```

### ✨ 核心特性

| 特性                | 描述                                                                 |
|---------------------|--------------------------------------------------------------------|
| **高效编码/解码**    | 支持标准URL编码和RFC 3986编码                                      |
| **精确解析**         | 支持灵活解析URL的各个组件                                          |
| **查询构建**         | 自动排序参数并生成标准查询字符串                                    |
| **并发安全**         | 所有方法线程安全                                                   |
| **错误处理**         | 提供详细的错误信息                                                |

### ⚠️ 注意事项
1. `Decode`方法会将"+"解码为空格
2. 使用`RawDecode`时需确保输入是RFC 3986编码
3. `ParseURL`的组件标志位使用位运算组合
4. 解析端口时返回字符串类型，需自行转换
5. 建议优先使用`DecodeTo`进行结构体绑定

### 🤝 参与贡献
[贡献指南](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [提交Issue](https://github.com/small-ek/antgo/issues)

---

## English

### 📖 Introduction

`aurl` is a high-performance URL processing library providing encoding, decoding, query building and parsing capabilities.  
Suitable for URL parameter handling, REST API development, web crawlers and other scenarios requiring precise URL manipulation.

GitHub URL: [github.com/small-ek/antgo/encoding/aurl](https://github.com/small-ek/antgo/encoding/aurl)

### 📦 Installation

```bash
go get github.com/small-ek/antgo/encoding/aurl
```

### 🚀 Quick Start

#### URL Encoding
```go
// Standard encoding
encoded := aurl.Encode("name=张三&age=30")

// RFC 3986 encoding
rawEncoded := aurl.RawEncode("a b~c")
```

#### URL Decoding
```go
decoded, _ := aurl.Decode("q%3Dgolang%26page%3D1")
rawDecoded, _ := aurl.RawDecode("a%20b%7Ec")
```

#### Query Building
```go
params := url.Values{
	"q":    []string{"golang"},
	"page": []string{"1"},
}
query := aurl.BuildQuery(params)
```

#### URL Parsing
```go
urlStr := "https://user:pass@example.com:8080/path?query=param#fragment"

// Parse all components
result, _ := aurl.ParseURL(urlStr, -1)

// Parse specific components
partial, _ := aurl.ParseURL(urlStr, 1|2|4) // scheme + host + port
```

### ✨ Key Features

| Feature             | Description                                                     |
|---------------------|-----------------------------------------------------------------|
| **Dual Encoding**   | Supports both standard and RFC 3986 encoding                   |
| **Precise Parsing** | Flexible component-based URL parsing                           |
| **Query Building**  | Auto-sorted parameter encoding                                 |
| **Concurrency Safe**| Thread-safe implementation                                     |
| **Error Handling**  | Detailed error messages                                        |

### ⚠️ Important Notes
1. `Decode` converts "+" to spaces
2. Ensure RFC 3986 compliance when using `RawDecode`
3. Component flags use bitwise combinations
4. Port numbers are returned as strings
5. Prefer struct binding with `DecodeTo`

### 🤝 Contributing
[Contribution Guide](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/small-ek/antgo/issues)

[⬆ Back to Top](#中文)