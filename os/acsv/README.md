# acsv - CSV 文件操作管理库 / CSV File Operation Manager

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`acsv` 是一个简单、灵活且高效的 CSV 文件操作管理库，旨在为 Go 项目提供一个简洁易用的 CSV 文件操作功能。它支持文件的读取、写入、更新、删除等操作，支持并发执行任务，保证数据安全。同时，`acsv` 提供了对 CSV 文件内容的直接管理，帮助用户轻松地对 CSV 文件进行增删改查（CRUD）操作。

GitHub 地址: [github.com/small-ek/antgo/os/acsv](https://github.com/small-ek/antgo/os/acsv)

### 📦 安装

```bash
go get github.com/small-ek/antgo/os/acsv
```

### 🚀 快速开始

#### 创建 CSV 实例
```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/os/acsv"
)

func main() {
	// 创建 CSV 实例
	csv, err := acsv.New("data.csv")
	if err != nil {
		fmt.Println("创建 CSV 实例失败:", err)
		return
	}

	// 输出当前 CSV 数据
	fmt.Println("当前 CSV 数据:", csv.Data)
}
```

#### 创建 CSV 文件并添加 BOM 头
```go
func main() {
	// 创建 CSV 实例
	csv, err := acsv.New("data.csv")
	if err != nil {
		fmt.Println("创建 CSV 实例失败:", err)
		return
	}

	// 创建 CSV 文件并添加 BOM 头
	err = csv.Create()
	if err != nil {
		fmt.Println("创建 CSV 文件失败:", err)
		return
	}

	// 输出 CSV 文件状态
	fmt.Println("CSV 文件创建成功，已添加 BOM 头")
}
```

#### 读取 CSV 文件内容
```go
func main() {
	// 创建 CSV 实例
	csv, err := acsv.New("data.csv")
	if err != nil {
		fmt.Println("创建 CSV 实例失败:", err)
		return
	}

	// 读取 CSV 文件内容
	err = csv.Read()
	if err != nil {
		fmt.Println("读取 CSV 文件失败:", err)
		return
	}

	// 输出 CSV 数据
	fmt.Println("CSV 数据:", csv.Data)
}
```

#### 写入 CSV 文件内容
```go
func main() {
	// 创建 CSV 实例
	csv, err := acsv.New("data.csv")
	if err != nil {
		fmt.Println("创建 CSV 实例失败:", err)
		return
	}

	// 添加一行数据
	csv.AddRow([]string{"Name", "Age", "Location"})

	// 写入 CSV 文件内容
	err = csv.Write()
	if err != nil {
		fmt.Println("写入 CSV 文件失败:", err)
		return
	}

	// 输出文件写入成功消息
	fmt.Println("CSV 文件已写入数据")
}
```

#### 删除 CSV 文件中的一行
```go
func main() {
	// 创建 CSV 实例
	csv, err := acsv.New("data.csv")
	if err != nil {
		fmt.Println("创建 CSV 实例失败:", err)
		return
	}

	// 读取 CSV 文件内容
	err = csv.Read()
	if err != nil {
		fmt.Println("读取 CSV 文件失败:", err)
		return
	}

	// 删除第一行数据
	err = csv.DeleteRow(0)
	if err != nil {
		fmt.Println("删除行失败:", err)
		return
	}

	// 写入更新后的 CSV 文件
	err = csv.Write()
	if err != nil {
		fmt.Println("写入 CSV 文件失败:", err)
		return
	}

	// 输出删除后的 CSV 数据
	fmt.Println("更新后的 CSV 数据:", csv.Data)
}
```

### ✨ 核心特性

| 特性                  | 描述                                                                 |
|-----------------------|----------------------------------------------------------------------|
| **秒级支持**           | 支持秒级 CSV 文件操作，满足精确处理需求                               |
| **多功能支持**         | 支持文件的创建、读取、写入、删除、更新等多种操作                     |
| **任务管理**           | 支持 CSV 数据的增删改查（CRUD）操作，并保证操作的线程安全             |
| **并发执行**           | 支持并发执行任务，提高任务执行效率                                   |
| **重试机制**           | 内置重试机制，自动重试失败的任务，提高任务成功率                       |
| **自动清理**           | 自动清理无效的任务，确保数据的一致性和稳定性                         |

### ⚠️ 注意事项

1. 确保 CSV 文件路径正确，避免因路径错误导致文件无法创建或读取。
2. 在高并发环境下使用时，注意 CSV 文件操作的线程安全性。
3. 对于大文件，建议采取批量写入操作以提高性能。
4. 对于长时间运行的任务，建议在任务执行中处理错误并进行日志记录。

### 🤝 参与贡献

[贡献指南](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [提交 Issue](https://github.com/small-ek/antgo/issues)

---

## English

### 📖 Introduction

`acsv` is a simple, flexible, and efficient CSV file operation management library designed to provide easy-to-use CSV file handling features for Go projects. It supports operations such as reading, writing, updating, and deleting CSV files, ensuring thread safety while supporting concurrent task execution. `acsv` helps you manage CSV file data directly, making CRUD operations on CSV files easy and efficient.

GitHub URL: [github.com/small-ek/antgo/os/acsv](https://github.com/small-ek/antgo/os/acsv)

### 📦 Installation

```bash
go get github.com/small-ek/antgo/os/acsv
```

### 🚀 Quick Start

#### Create a CSV Instance
```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/os/acsv"
)

func main() {
	// Create CSV instance
	csv, err := acsv.New("data.csv")
	if err != nil {
		fmt.Println("Failed to create CSV instance:", err)
		return
	}

	// Output current CSV data
	fmt.Println("Current CSV Data:", csv.Data)
}
```

#### Create a CSV File and Add BOM Header
```go
func main() {
	// Create CSV instance
	csv, err := acsv.New("data.csv")
	if err != nil {
		fmt.Println("Failed to create CSV instance:", err)
		return
	}

	// Create CSV file and add BOM header
	err = csv.Create()
	if err != nil {
		fmt.Println("Failed to create CSV file:", err)
		return
	}

	// Output CSV file status
	fmt.Println("CSV file created with BOM header")
}
```

#### Read CSV File Contents
```go
func main() {
	// Create CSV instance
	csv, err := acsv.New("data.csv")
	if err != nil {
		fmt.Println("Failed to create CSV instance:", err)
		return
	}

	// Read CSV file contents
	err = csv.Read()
	if err != nil {
		fmt.Println("Failed to read CSV file:", err)
		return
	}

	// Output CSV data
	fmt.Println("CSV Data:", csv.Data)
}
```

#### Write to CSV File
```go
func main() {
	// Create CSV instance
	csv, err := acsv.New("data.csv")
	if err != nil {
		fmt.Println("Failed to create CSV instance:", err)
		return
	}

	// Add a row
	csv.AddRow([]string{"Name", "Age", "Location"})

	// Write to CSV file
	err = csv.Write()
	if err != nil {
		fmt.Println("Failed to write to CSV file:", err)
		return
	}

	// Output success message
	fmt.Println("CSV file written successfully")
}
```

#### Delete a Row from CSV File
```go
func main() {
	// Create CSV instance
	csv, err := acsv.New("data.csv")
	if err != nil {
		fmt.Println("Failed to create CSV instance:", err)
		return
	}

	// Read CSV file contents
	err = csv.Read()
	if err != nil {
		fmt.Println("Failed to read CSV file:", err)
		return
	}

	// Delete the first row
	err = csv.DeleteRow(0)
	if err != nil {
		fmt.Println("Failed to delete row:", err)
		return
	}

	// Write updated CSV file
	err = csv.Write()
	if err != nil {
		fmt.Println("Failed to write CSV file:", err)
		return
	}

	// Output updated CSV data
	fmt.Println("Updated CSV Data:", csv.Data)
}
```

### ✨ Key Features

| Feature               | Description                                                           |
|-----------------------|-----------------------------------------------------------------------|
| **Second-Level Support** | Supports second-level CSV file operations for precise handling         |
| **Multi-Method Support** | Supports file creation, reading, writing, deleting, updating, and more |
| **Task Management**     | Easy CRUD operations on CSV data with thread safety                     |
| **Concurrent Execution** | Supports concurrent execution of tasks for improved performance         |
| **Retry Mechanism**     | Built-in retry mechanism to improve

task success rates                  |
| **Auto Cleanup**        | Automatically cleans up invalid tasks to ensure data integrity and stability |

### ⚠️ Important Notes

1. Ensure the CSV file path is correct to avoid issues when creating or reading files.
2. In high-concurrency environments, be mindful of thread safety when performing file operations.
3. For large files, consider batch writing to improve performance.
4. For long-running tasks, ensure proper error handling and logging.

### 🤝 Contributing

[Contribution Guide](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/small-ek/antgo/issues)

[⬆ Back to Top](#中文)
