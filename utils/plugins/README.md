# plugins - 插件管理工具库 / Plugin Management Utility

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`plugins` 是一个简单易用的插件管理工具库，提供插件的注册、卸载和列表查看功能。  
该工具库适用于需要动态加载和管理插件的场景，例如插件化架构的系统或模块化应用。  
通过 `PluginManager` 结构体，您可以轻松管理插件的生命周期，并确保线程安全。

GitHub 地址: [github.com/small-ek/antgo/utils/plugins](https://github.com/small-ek/antgo/utils/plugins)

### 📦 安装

使用 `go get` 命令进行安装：

```bash
go get github.com/small-ek/antgo/utils/plugins
```

### 🚀 快速开始

#### 初始化插件管理器

使用 `New` 方法初始化插件管理器的单例实例。  
示例代码：

```go
package main

import (
	"github.com/small-ek/antgo/utils/plugins"
)

func main() {
	manager := plugins.New()
	fmt.Println("插件管理器初始化成功！")
}
```

#### 注册插件

使用 `Register` 方法将插件注册到管理器中。  
示例代码：

```go
package main

import (
	"fmt"
	"log"

	"github.com/small-ek/antgo/utils/plugins"
)

type MyPlugin struct{}

func (p *MyPlugin) Before() interface{} {
	fmt.Println("执行 Before 方法")
	return nil
}

func (p *MyPlugin) After(data ...interface{}) interface{} {
	fmt.Println("执行 After 方法")
	return nil
}

func main() {
	manager := plugins.New()
	err := manager.Register("myPlugin", &MyPlugin{})
	if err != nil {
		log.Fatalf("插件注册失败: %v", err)
	}
	fmt.Println("插件注册成功！")
}
```

#### 卸载插件

使用 `Uninstall` 方法从管理器中卸载插件。  
示例代码：

```go
package main

import (
	"fmt"
	"log"

	"github.com/small-ek/antgo/utils/plugins"
)

func main() {
	manager := plugins.New()
	err := manager.Uninstall("myPlugin")
	if err != nil {
		log.Fatalf("插件卸载失败: %v", err)
	}
	fmt.Println("插件卸载成功！")
}
```

#### 查看插件列表

使用 `List` 方法获取所有已注册插件的列表。  
示例代码：

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/plugins"
)

func main() {
	manager := plugins.New()
	pluginList := plugins.List()
	fmt.Println("已注册插件列表:", pluginList)
}
```

### ✨ 核心特性

| 特性                | 描述                                                      |
|---------------------|-----------------------------------------------------------|
| **插件注册**        | 支持动态注册插件，确保插件名称唯一                         |
| **插件卸载**        | 支持按名称卸载插件，释放资源                               |
| **插件列表**        | 提供查看所有已注册插件的功能                               |
| **线程安全**        | 使用读写锁 (`sync.RWMutex`) 确保并发安全                   |
| **简单易用**        | 提供直观的 API，方便在各类项目中集成和使用                  |

### ⚠️ 注意事项

1. 插件名称必须唯一，重复注册同名插件会导致错误。
2. 卸载插件时需确保插件名称正确，否则会返回错误。
3. 插件管理器为单例模式，全局共享同一个实例。

### 🤝 参与贡献

- [贡献指南](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md)
- [提交 Issue](https://github.com/small-ek/antgo/issues)

[⬆ 返回顶部](#plugins---插件管理工具库--plugin-management-utility)

---

## English

### 📖 Introduction

`plugins` is a straightforward and easy-to-use plugin management utility library that provides functions for registering, uninstalling, and listing plugins.  
It is well-suited for scenarios requiring dynamic loading and management of plugins, such as plugin-based architectures or modular applications.  
Through the `PluginManager` struct, you can easily manage the lifecycle of plugins and ensure thread safety.

GitHub URL: [github.com/small-ek/antgo/utils/plugins](https://github.com/small-ek/antgo/utils/plugins)

### 📦 Installation

Install via `go get`:

```bash
go get github.com/small-ek/antgo/utils/plugins
```

### 🚀 Quick Start

#### Initializing the Plugin Manager

Use the `New` function to initialize the singleton instance of the plugin manager.  
Example:

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/plugins"
)

func main() {
	manager := plugins.New()
	fmt.Println("Plugin manager initialized successfully!")
}
```

#### Registering a Plugin

Use the `Register` function to add a plugin to the manager.  
Example:

```go
package main

import (
	"fmt"
	"log"

	"github.com/small-ek/antgo/utils/plugins"
)

type MyPlugin struct{}

func (p *MyPlugin) Before() interface{} {
	fmt.Println("Executing Before method")
	return nil
}

func (p *MyPlugin) After(data ...interface{}) interface{} {
	fmt.Println("Executing After method")
	return nil
}

func main() {
	manager := plugins.New()
	err := manager.Register("myPlugin", &MyPlugin{})
	if err != nil {
		log.Fatalf("Failed to register plugin: %v", err)
	}
	fmt.Println("Plugin registered successfully!")
}
```

#### Uninstalling a Plugin

Use the `Uninstall` function to remove a plugin from the manager.  
Example:

```go
package main

import (
	"fmt"
	"log"

	"github.com/small-ek/antgo/utils/plugins"
)

func main() {
	manager := plugins.New()
	err := manager.Uninstall("myPlugin")
	if err != nil {
		log.Fatalf("Failed to uninstall plugin: %v", err)
	}
	fmt.Println("Plugin uninstalled successfully!")
}
```

#### Listing Plugins

Use the `List` function to retrieve a list of all registered plugins.  
Example:

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/plugins"
)

func main() {
	manager := plugins.New()
	pluginList := plugins.List()
	fmt.Println("Registered plugins:", pluginList)
}
```

### ✨ Key Features

| Feature               | Description                                                      |
|-----------------------|------------------------------------------------------------------|
| **Plugin Registration** | Supports dynamic plugin registration with unique names           |
| **Plugin Uninstallation** | Allows uninstalling plugins by name, freeing resources           |
| **Plugin Listing**    | Provides a list of all registered plugins                        |
| **Thread Safety**     | Ensures concurrency safety using `sync.RWMutex`                  |
| **Simplicity**        | Offers an intuitive API for easy integration into projects       |

### ⚠️ Important Notes

1. Plugin names must be unique; registering a duplicate name will result in an error.
2. Ensure the correct plugin name is used when uninstalling, or an error will be returned.
3. The plugin manager is a singleton, sharing a single instance globally.

### 🤝 Contributing

- [Contribution Guide](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md)
- [Open an Issue](https://github.com/small-ek/antgo/issues)

[⬆ Back to Top](#plugins---plugin-management-utility)