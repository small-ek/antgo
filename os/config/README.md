# config - 配置管理库 / Configuration Management Library

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`config` 模块是一个基于 [viper](https://github.com/spf13/viper) 的配置管理库，专为 Go 项目设计。它不仅提供了本地配置文件的加载与管理，还支持远程配置（例如 ETCD 和其他远程配置提供者）的读取和监听。该模块封装了常用的配置获取方法，帮助开发者更方便地在项目中管理配置参数。

GitHub 地址: [github.com/small-ek/antgo/os/config](https://github.com/small-ek/antgo/os/config)

### 📦 安装

```bash
go get github.com/small-ek/antgo/os/config
```

### 🚀 快速开始

#### 加载本地配置文件

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/os/config"
)

func main() {
	// 传入本地配置文件路径，支持 toml、json、yaml 等格式
	cfg := config.New("config.toml")
	if err := cfg.Register(); err != nil {
		panic(err)
	}

	// 获取配置项
	port := config.GetInt("port")
	name := config.GetString("name")
	fmt.Printf("应用启动于端口：%d, 名称：%s\n", port, name)
}
```

#### 合并多个配置文件

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/os/config"
)

func main() {
	// 先加载基础配置文件
	cfg := config.New("base_config.toml")
	if err := cfg.Register(); err != nil {
		panic(err)
	}

	// 通过 AddConfigFile 方法合并新的配置文件，支持动态监听文件变更
	if err := config.AddConfigFile("override_config.toml"); err != nil {
		panic(err)
	}

	// 获取合并后的配置项
	debug := config.GetBool("debug")
	fmt.Printf("调试模式：%t\n", debug)
}
```

#### 使用全局辅助函数设置和获取配置

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/os/config"
)

func main() {
	// 初始化空配置（配置文件内容为空也可）
	cfg := config.New("empty_config.toml")
	if err := cfg.Register(); err != nil {
		panic(err)
	}

	// 直接设置配置键值对
	config.SetKey("custom.setting", "customValue")

	// 获取设置的配置项
	value := config.GetString("custom.setting")
	fmt.Printf("自定义配置：%s\n", value)
}
```

#### 远程配置（ETCD3）示例

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/os/config"
)

func main() {
	// 初始化本地配置（此处仅用于演示，实际使用时可根据需要传入其他参数）
	cfg := config.New()
	
	// 连接 ETCD3 并加载配置
	// hosts: ETCD3 服务器地址列表
	// paths: ETCD 中存储配置的键（支持文件名格式，例如 "config.toml"）
	// username 和 pwd 为 ETCD3 的认证信息
	err := cfg.Etcd3([]string{"127.0.0.1:2379"}, []string{"config.toml"}, "user", "password")
	if err != nil {
		panic(err)
	}

	// 读取 ETCD3 中的配置项
	appName := config.GetString("app.name")
	fmt.Printf("应用名称：%s\n", appName)

	// 程序将持续监听 ETCD3 配置变化
	select {}
}
```

### ✨ 核心特性

| 特性                  | 描述                                                                 |
|-----------------------|----------------------------------------------------------------------|
| **多格式支持**         | 支持 TOML、JSON、YAML 等多种配置文件格式                              |
| **远程配置加载**       | 支持 ETCD3 以及其他远程配置服务的连接和配置动态更新                     |
| **配置合并**           | 通过合并多个配置文件，实现配置覆盖与扩展                                |
| **全局辅助函数**       | 提供便捷的 Get 和 Set 系列函数，方便全局配置的读取与设置                 |
| **动态监听**           | 支持文件变更和远程配置更新时自动重新加载配置，保持配置最新                 |
| **线程安全**           | 支持并发读取与更新，适用于高并发场景                                     |

### ⚠️ 注意事项

1. 确保传入的配置文件路径和文件格式正确，否则可能会导致读取失败。
2. 远程配置部分需要确保网络连通性及正确的认证信息。
3. 在使用配置合并和监听功能时，请注意性能与线程安全问题。
4. 根据项目需求合理设计配置的层级与覆盖规则。

### 🤝 参与贡献

[贡献指南](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [提交 Issue](https://github.com/small-ek/antgo/issues)

---

## English

### 📖 Introduction

The `config` package is a configuration management library designed for Go projects and built on top of [viper](https://github.com/spf13/viper). It supports loading and managing local configuration files as well as remote configurations (such as from ETCD and other remote providers). The package also provides a set of helper functions to conveniently retrieve and set configuration parameters in your project.

GitHub URL: [github.com/small-ek/antgo/os/config](https://github.com/small-ek/antgo/os/config)

### 📦 Installation

```bash
go get github.com/small-ek/antgo/os/config
```

### 🚀 Quick Start

#### Loading Local Configuration Files

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/os/config"
)

func main() {
	// Initialize configuration by providing the local configuration file path.
	// Supported formats include TOML, JSON, YAML, etc.
	cfg := config.New("config.toml")
	if err := cfg.Register(); err != nil {
		panic(err)
	}

	// Retrieve configuration values
	port := config.GetInt("port")
	name := config.GetString("name")
	fmt.Printf("Application is running on port: %d, Name: %s\n", port, name)
}
```

#### Merging Multiple Configuration Files

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/os/config"
)

func main() {
	// Load the base configuration file
	cfg := config.New("base_config.toml")
	if err := cfg.Register(); err != nil {
		panic(err)
	}

	// Merge an additional configuration file.
	// The merged file supports dynamic watching for file changes.
	if err := config.AddConfigFile("override_config.toml"); err != nil {
		panic(err)
	}

	// Retrieve merged configuration values
	debug := config.GetBool("debug")
	fmt.Printf("Debug mode: %t\n", debug)
}
```

#### Using Global Helper Functions

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/os/config"
)

func main() {
	// Initialize an empty configuration (even empty files can be registered).
	cfg := config.New("empty_config.toml")
	if err := cfg.Register(); err != nil {
		panic(err)
	}

	// Directly set a configuration key-value pair.
	config.SetKey("custom.setting", "customValue")

	// Retrieve the custom configuration value.
	value := config.GetString("custom.setting")
	fmt.Printf("Custom configuration: %s\n", value)
}
```

#### Remote Configuration (ETCD3) Example

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/os/config"
)

func main() {
	// Initialize local configuration (parameters can be adjusted as needed)
	cfg := config.New()
	
	// Connect to ETCD3 and load configurations.
	// hosts: list of ETCD3 server addresses.
	// paths: keys in ETCD storing the configuration (supports file-like naming such as "config.toml").
	// username and pwd are for ETCD3 authentication.
	err := cfg.Etcd3([]string{"127.0.0.1:2379"}, []string{"config.toml"}, "user", "password")
	if err != nil {
		panic(err)
	}

	// Retrieve configuration values from ETCD3
	appName := config.GetString("app.name")
	fmt.Printf("Application Name: %s\n", appName)

	// The application will continuously listen for configuration changes from ETCD3.
	select {}
}
```

### ✨ Key Features

| Feature                     | Description                                                                   |
|-----------------------------|-------------------------------------------------------------------------------|
| **Multi-Format Support**    | Supports configuration files in TOML, JSON, YAML, etc.                        |
| **Remote Configuration**    | Connects to ETCD3 and other remote providers to load and dynamically update configurations |
| **Configuration Merging**   | Merge multiple configuration files for overriding and extending settings      |
| **Global Helper Functions** | Provides easy-to-use Get and Set functions for configuration management         |
| **Dynamic Watching**        | Automatically reloads configuration on file changes or remote updates           |
| **Thread Safety**           | Safe for concurrent access in high-concurrency environments                     |

### ⚠️ Important Notes

1. Ensure that the configuration file paths and formats provided are correct to avoid read failures.
2. For remote configurations, verify network connectivity and proper authentication.
3. When using configuration merging and watching features, consider performance and thread safety.
4. Design the configuration hierarchy and override rules based on your project requirements.

### 🤝 Contributing

[Contribution Guide](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/small-ek/antgo/issues)

[⬆ Back to Top](#中文)