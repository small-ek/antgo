# antgo/encoding/aini - INI 配置文件处理库 / INI Configuration Processing Library

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`antgo/encoding/aini` 是一个高效的INI配置文件处理库，支持INI文件的解析、编码以及转换为JSON格式。  
适用于配置文件读取、写入、格式转换等场景。

GitHub地址: [github.com/small-ek/antgo/encoding/aini](https://github.com/small-ek/antgo/encoding/aini)

### 📦 安装

```bash
go get github.com/small-ek/antgo/encoding/aini
```

### 🚀 快速开始

#### 基本用法
```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/encoding/aini"
)

func main() {
	// 示例INI内容
	iniContent := `
[database]
host = 127.0.0.1
port = 3306
user = root
password = secret

[server]
host = 0.0.0.0
port = 8080
`

	// 解析INI内容
	iniMap, err := aini.Decode([]byte(iniContent))
	if err != nil {
		fmt.Println("解析失败:", err)
		return
	}
	fmt.Println(iniMap) // 输出: map[database:map[host:127.0.0.1 port:3306 user:root password:secret] server:map[host:0.0.0.0 port:8080]]

	// 将解析后的数据编码为INI格式
	encodedINI, err := aini.Encode(iniMap)
	if err != nil {
		fmt.Println("编码失败:", err)
		return
	}
	fmt.Println(string(encodedINI)) // 输出: [database]\nhost = 127.0.0.1\nport = 3306\nuser = root\npassword = secret\n\n[server]\nhost = 0.0.0.0\nport = 8080\n
}
```

#### 转换为JSON
```go
func main() {
	iniContent := `
[database]
host = 127.0.0.1
port = 3306
`

	// 将INI内容转换为JSON
	jsonData, err := aini.ToJson([]byte(iniContent))
	if err != nil {
		fmt.Println("转换失败:", err)
		return
	}
	fmt.Println(string(jsonData)) // 输出: {"database":{"host":"127.0.0.1","port":"3306"}}
}
```

### ✨ 核心特性

| 特性                | 描述                                                                 |
|---------------------|--------------------------------------------------------------------|
| **高效解析**         | 快速解析INI文件内容，支持嵌套节和键值对                             |
| **编码支持**         | 将解析后的数据重新编码为INI格式                                     |
| **JSON转换**         | 支持将INI内容转换为JSON格式，便于进一步处理                         |
| **并发安全**         | 所有函数线程安全，支持高并发场景                                   |

### ⚠️ 注意事项
1. INI文件必须包含至少一个有效的节（section），否则会返回错误。
2. 键值对的分隔符为 `=`，且键和值两端的空格会被自动去除。
3. 注释以 `;` 或 `#` 开头，解析时会自动忽略。
4. 返回结果均为新副本，原始数据不会被修改。

### 🤝 参与贡献
[贡献指南](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [提交Issue](https://github.com/small-ek/antgo/issues)

---

## English

### 📖 Introduction

`antgo/encoding/aini` is a high-performance INI configuration file processing library, supporting parsing, encoding, and conversion to JSON format.  
Ideal for reading, writing, and converting configuration files.

GitHub URL: [github.com/small-ek/antgo/encoding/aini](https://github.com/small-ek/antgo/encoding/aini)

### 📦 Installation

```bash
go get github.com/small-ek/antgo/encoding/aini
```

### 🚀 Quick Start

#### Basic Usage
```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/encoding/aini"
)

func main() {
	// Example INI content
	iniContent := `
[database]
host = 127.0.0.1
port = 3306
user = root
password = secret

[server]
host = 0.0.0.0
port = 8080
`

	// Parse INI content
	iniMap, err := aini.Decode([]byte(iniContent))
	if err != nil {
		fmt.Println("Parse failed:", err)
		return
	}
	fmt.Println(iniMap) // Output: map[database:map[host:127.0.0.1 port:3306 user:root password:secret] server:map[host:0.0.0.0 port:8080]]

	// Encode parsed data back to INI format
	encodedINI, err := aini.Encode(iniMap)
	if err != nil {
		fmt.Println("Encode failed:", err)
		return
	}
	fmt.Println(string(encodedINI)) // Output: [database]\nhost = 127.0.0.1\nport = 3306\nuser = root\npassword = secret\n\n[server]\nhost = 0.0.0.0\nport = 8080\n
}
```

#### Convert to JSON
```go
func main() {
	iniContent := `
[database]
host = 127.0.0.1
port = 3306
`

	// Convert INI content to JSON
	jsonData, err := aini.ToJson([]byte(iniContent))
	if err != nil {
		fmt.Println("Conversion failed:", err)
		return
	}
	fmt.Println(string(jsonData)) // Output: {"database":{"host":"127.0.0.1","port":"3306"}}
}
```

### ✨ Key Features

| Feature             | Description                                                        |
|---------------------|--------------------------------------------------------------------|
| **Efficient Parsing**| Fast parsing of INI files with support for nested sections and key-value pairs |
| **Encoding Support** | Encode parsed data back to INI format                             |
| **JSON Conversion**  | Convert INI content to JSON format for further processing         |
| **Thread-Safe**     | All functions are concurrency-ready                               |

### ⚠️ Important Notes
1. INI files must contain at least one valid section, otherwise an error will be returned.
2. Key-value pairs are separated by `=`, and spaces around keys and values are automatically trimmed.
3. Comments starting with `;` or `#` are ignored during parsing.
4. All results are new copies, original data remains unchanged.

### 🤝 Contributing
[Contribution Guide](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/small-ek/antgo/issues)

[⬆ Back to Top](#中文)
