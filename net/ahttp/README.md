# ahttp - HTTP客户端库 / HTTP Client Utilities

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`ahttp` 是一个高效、易用的HTTP客户端库，提供了简洁的API来发送HTTP请求并处理响应。它支持GET、POST、PUT、DELETE等多种HTTP方法，并且可以轻松地设置请求头、请求参数、超时时间等。`ahttp` 还内置了重试机制、连接池管理等功能，适用于各种HTTP请求场景。

GitHub地址: [github.com/small-ek/antgo/net/ahttp](https://github.com/small-ek/antgo/net/ahttp)

### 📦 安装

```bash
go get github.com/small-ek/antgo/net/ahttp
```

### 🚀 快速开始

#### 发送GET请求
```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/net/ahttp"
)

func main() {
	// 创建一个HTTP客户端
	client := ahttp.New()

	// 发送GET请求
	response, err := client.Get("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}

	// 打印响应内容
	fmt.Println("响应状态码:", response.StatusCode())
	fmt.Println("响应体:", response.String())
}
```

#### 发送POST请求
```go
func main() {
	client := ahttp.New()

	// 设置请求体
	body := map[string]interface{}{
		"title":  "foo",
		"body":   "bar",
		"userId": 1,
	}

	// 发送POST请求
	response, err := client.Post("https://jsonplaceholder.typicode.com/posts", body)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}

	// 打印响应内容
	fmt.Println("响应状态码:", response.StatusCode())
	fmt.Println("响应体:", response.String())
}
```

### 🔧 高级用法

#### 设置请求头
```go
func main() {
	client := ahttp.New()

	// 设置自定义请求头
	client.SetHeader("Authorization", "Bearer token123")
	client.SetHeader("Content-Type", "application/json")

	response, err := client.Get("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}

	fmt.Println("响应体:", response.String())
}
```

#### 设置超时时间
```go
func main() {
	client := ahttp.New()

	// 设置请求超时时间为5秒
	client.SetTimeout(5 * time.Second)

	response, err := client.Get("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}

	fmt.Println("响应体:", response.String())
}
```

#### 使用重试机制
```go
func main() {
	client := ahttp.New()

	// 设置重试次数为3次
	client.SetRetryCount(3)

	response, err := client.Get("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}

	fmt.Println("响应体:", response.String())
}
```

### ✨ 核心特性

| 特性                | 描述                                                                 |
|---------------------|--------------------------------------------------------------------|
| **多方法支持**       | 支持GET、POST、PUT、DELETE等多种HTTP方法                            |
| **请求头设置**       | 轻松设置自定义请求头                                               |
| **超时控制**         | 支持设置请求超时时间                                               |
| **重试机制**         | 内置请求重试机制，提高请求成功率                                   |
| **连接池管理**       | 自动管理HTTP连接池，提升性能                                       |

### ⚠️ 注意事项
1. 确保请求的URL是有效的
2. 设置合理的超时时间以避免请求长时间挂起
3. 使用重试机制时，注意服务器的负载情况
4. 对于敏感数据，建议使用HTTPS协议

### 🤝 参与贡献
[贡献指南](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [提交Issue](https://github.com/small-ek/antgo/issues)

---

## English

### 📖 Introduction

`ahttp` is an efficient and easy-to-use HTTP client library that provides a simple API for sending HTTP requests and handling responses. It supports various HTTP methods such as GET, POST, PUT, DELETE, and allows easy configuration of request headers, parameters, and timeout settings. `ahttp` also includes features like retry mechanism and connection pool management, making it suitable for various HTTP request scenarios.

GitHub URL: [github.com/small-ek/antgo/net/ahttp](https://github.com/small-ek/antgo/net/ahttp)

### 📦 Installation

```bash
go get github.com/small-ek/antgo/net/ahttp
```

### 🚀 Quick Start

#### Sending a GET Request
```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/net/ahttp"
)

func main() {
	client := ahttp.New()

	response, err := client.Get("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}

	fmt.Println("Response status code:", response.StatusCode())
	fmt.Println("Response body:", response.String())
}
```

#### Sending a POST Request
```go
func main() {
	client := ahttp.New()

	body := map[string]interface{}{
		"title":  "foo",
		"body":   "bar",
		"userId": 1,
	}

	response, err := client.Post("https://jsonplaceholder.typicode.com/posts", body)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}

	fmt.Println("Response status code:", response.StatusCode())
	fmt.Println("Response body:", response.String())
}
```

### 🔧 Advanced Usage

#### Setting Request Headers
```go
func main() {
	client := ahttp.New()

	client.SetHeader("Authorization", "Bearer token123")
	client.SetHeader("Content-Type", "application/json")

	response, err := client.Get("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}

	fmt.Println("Response body:", response.String())
}
```

#### Setting Timeout
```go
func main() {
	client := ahttp.New()

	client.SetTimeout(5 * time.Second)

	response, err := client.Get("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}

	fmt.Println("Response body:", response.String())
}
```

#### Using Retry Mechanism
```go
func main() {
	client := ahttp.New()

	client.SetRetryCount(3)

	response, err := client.Get("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}

	fmt.Println("Response body:", response.String())
}
```

### ✨ Key Features

| Feature             | Description                                                     |
|---------------------|-----------------------------------------------------------------|
| **Multi-method**    | Supports GET, POST, PUT, DELETE, and more                       |
| **Header Setting**  | Easy configuration of custom request headers                   |
| **Timeout Control** | Supports setting request timeout                               |
| **Retry Mechanism** | Built-in retry mechanism to improve request success rate       |
| **Connection Pool** | Automatic management of HTTP connection pool for better performance |

### ⚠️ Important Notes
1. Ensure the request URL is valid
2. Set a reasonable timeout to avoid long hanging requests
3. Be mindful of server load when using the retry mechanism
4. Use HTTPS for sensitive data

### 🤝 Contributing
[Contribution Guide](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/small-ek/antgo/issues)

[⬆ Back to Top](#中文)