# acron - Cron 任务调度管理库 / Cron Task Scheduler Manager

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`acron` 是一个简单、灵活且高效的 Cron 任务调度管理库，旨在为 Go 项目提供一个简洁易用的 Cron 定时任务管理功能。它支持秒级任务调度，可以轻松添加、删除、查询任务，管理任务ID，并且支持并发执行任务。`acron` 还内置了任务清理和重试机制，适用于各种高并发和高可用的任务调度场景。

GitHub 地址: [github.com/small-ek/antgo/os/acron](https://github.com/small-ek/antgo/os/acron)

### 📦 安装

```bash
go get github.com/small-ek/antgo/os/acron
```

### 🚀 快速开始

#### 创建 Cron 实例
```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/os/acron"
)

func main() {
	// 创建 Cron 管理实例
	crontab := acron.New()

	// 添加一个定时任务
	err := crontab.AddByID("task1", "* * * * *", cron.FuncJob(func() {
		fmt.Println("任务1执行")
	}))
	if err != nil {
		fmt.Println("添加任务失败:", err)
		return
	}

	// 启动 Cron 引擎
	crontab.Start()

	// 输出当前有效的任务ID
	fmt.Println("有效任务ID:", crontab.IDs())

	// 停止 Cron 引擎
	crontab.Stop()
}
```

#### 添加 Cron 任务
```go
func main() {
	crontab := acron.New()

	// 添加任务
	err := crontab.AddByID("task2", "*/5 * * * *", cron.FuncJob(func() {
		fmt.Println("任务2执行，每5分钟一次")
	}))
	if err != nil {
		fmt.Println("添加任务失败:", err)
		return
	}

	// 启动 Cron 引擎
	crontab.Start()
}
```

### 🔧 高级用法

#### 设置任务执行的函数
```go
func main() {
	crontab := acron.New()

	// 设置定时任务的执行函数
	err := crontab.AddByFunc("task3", "*/10 * * * *", func() {
		fmt.Println("任务3执行，每10分钟一次")
	})
	if err != nil {
		fmt.Println("添加任务失败:", err)
		return
	}

	// 启动 Cron 引擎
	crontab.Start()
}
```

#### 删除 Cron 任务
```go
func main() {
	crontab := acron.New()

	// 添加任务
	crontab.AddByID("task4", "* * * * *", cron.FuncJob(func() {
		fmt.Println("任务4执行")
	}))

	// 删除任务
	crontab.DelByID("task4")

	// 启动 Cron 引擎
	crontab.Start()

	// 输出当前有效的任务ID
	fmt.Println("当前有效的任务ID:", crontab.IDs())
}
```

### ✨ 核心特性

| 特性                  | 描述                                                                 |
|-----------------------|----------------------------------------------------------------------|
| **秒级支持**           | 支持秒级Cron表达式，满足各种高精度定时任务需求                         |
| **多方法支持**         | 支持添加任务、删除任务、查询任务等多种操作                            |
| **任务管理**           | 轻松管理任务ID，支持任务的添加、删除、查询、任务清理等操作             |
| **并发执行**           | 支持并发执行任务，提高任务的执行效率                                  |
| **重试机制**           | 内置重试机制，自动重试失败的任务，提高任务成功率                       |
| **任务清理**           | 自动清理无效的任务，确保任务调度的稳定性                               |

### ⚠️ 注意事项

1. 确保任务ID是唯一的，避免重复添加同一任务。
2. 设置合理的任务时间间隔，避免过高频率的任务执行。
3. 在高并发环境下使用时，确保考虑到任务执行的并发性。
4. 对于长时间运行的任务，建议在任务执行中处理错误并进行日志记录。

### 🤝 参与贡献

[贡献指南](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [提交 Issue](https://github.com/small-ek/antgo/issues)

---

## English

### 📖 Introduction

`acron` is a simple, flexible, and efficient Cron task scheduling manager library designed to provide easy-to-use Cron scheduling features for Go projects. It supports second-level task scheduling, allowing you to easily add, remove, and query tasks, manage task IDs, and handle task execution concurrently. `acron` also includes features such as task cleanup and retry mechanism, making it suitable for high-concurrency and high-availability task scheduling scenarios.

GitHub URL: [github.com/small-ek/antgo/os/acron](https://github.com/small-ek/antgo/os/acron)

### 📦 Installation

```bash
go get github.com/small-ek/antgo/os/acron
```

### 🚀 Quick Start

#### Create a Cron Instance
```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/os/acron"
)

func main() {
	// Create Cron manager instance
	crontab := acron.New()

	// Add a task
	err := crontab.AddByID("task1", "* * * * *", cron.FuncJob(func() {
		fmt.Println("Task 1 executed")
	}))
	if err != nil {
		fmt.Println("Failed to add task:", err)
		return
	}

	// Start the Cron engine
	crontab.Start()

	// Print current valid task IDs
	fmt.Println("Valid Task IDs:", crontab.IDs())

	// Stop the Cron engine
	crontab.Stop()
}
```

#### Adding Cron Tasks
```go
func main() {
	crontab := acron.New()

	// Add a task
	err := crontab.AddByID("task2", "*/5 * * * *", cron.FuncJob(func() {
		fmt.Println("Task 2 executed every 5 minutes")
	}))
	if err != nil {
		fmt.Println("Failed to add task:", err)
		return
	}

	// Start the Cron engine
	crontab.Start()
}
```

### 🔧 Advanced Usage

#### Set Task Execution Function
```go
func main() {
	crontab := acron.New()

	// Set the function for the scheduled task
	err := crontab.AddByFunc("task3", "*/10 * * * *", func() {
		fmt.Println("Task 3 executed every 10 minutes")
	})
	if err != nil {
		fmt.Println("Failed to add task:", err)
		return
	}

	// Start the Cron engine
	crontab.Start()
}
```

#### Deleting Cron Tasks
```go
func main() {
	crontab := acron.New()

	// Add a task
	crontab.AddByID("task4", "* * * * *", cron.FuncJob(func() {
		fmt.Println("Task 4 executed")
	}))

	// Delete a task
	crontab.DelByID("task4")

	// Start the Cron engine
	crontab.Start()

	// Print current valid task IDs
	fmt.Println("Current valid task IDs:", crontab.IDs())
}
```

### ✨ Key Features

| Feature               | Description                                                           |
|-----------------------|-----------------------------------------------------------------------|
| **Second-Level Support** | Supports second-level Cron expressions for high-precision scheduling |
| **Multi-method Support** | Supports adding, deleting, querying tasks, and more                  |
| **Task Management**     | Easily manage task IDs, support task addition, deletion, query, and cleanup |
| **Concurrent Execution** | Supports concurrent task execution for improved performance           |
| **Retry Mechanism**     | Built-in retry mechanism for improving task success rate              |
| **Task Cleanup**        | Automatically cleans up invalid tasks to ensure task scheduler stability |

### ⚠️ Important Notes

1. Ensure that task IDs are unique to avoid adding the same task repeatedly.
2. Set reasonable task intervals to avoid tasks running too frequently.
3. Consider concurrency when using it in high-concurrency environments.
4. For long-running tasks, ensure proper error handling and logging.

### 🤝 Contributing

[Contribution Guide](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/small-ek/antgo/issues)

[⬆ Back to Top](#中文)