# password - 安全密码处理工具库 / Secure Password Utility

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`password` 是一个简单易用的密码处理工具库，提供基于 bcrypt 算法的密码哈希生成和验证功能。  
该工具库适用于用户密码存储、认证系统以及任何需要安全密码处理的场景。  
每次生成的哈希值都包含随机盐、版本信息和配置参数，确保即使相同密码多次生成也会有不同的哈希值。

GitHub 地址: [github.com/small-ek/antgo/utils/password](https://github.com/small-ek/antgo/utils/password)

### 📦 安装

使用 `go get` 命令进行安装：

```bash
go get github.com/small-ek/antgo/utils/password
```

### 🚀 快速开始

#### 生成密码哈希

使用 `Generate` 方法可以根据原始密码生成安全的 bcrypt 哈希。  
示例代码：

```go
package main

import (
	"fmt"
	"log"

	"github.com/small-ek/antgo/utils/password"
)

func main() {
	rawPassword := "mySecureP@ssw0rd"
	hashedPassword, err := password.Generate(rawPassword)
	if err != nil {
		log.Fatalf("生成密码哈希失败: %v", err)
	}
	fmt.Println("生成的哈希值:", hashedPassword)
}
```

#### 验证密码

使用 `Verify` 方法来验证用户输入的密码是否与存储的哈希值匹配。  
示例代码：

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/password"
)

func main() {
	storedHash := "$2a$08$u5hX9iIwBzr7rN7w5EjZne.Q/XO6zC8Q.f56kBvEUp6kT8YHvXGa6" // 示例哈希值
	inputPassword := "mySecureP@ssw0rd"

	if password.Verify(storedHash, inputPassword) {
		fmt.Println("密码验证成功！")
	} else {
		fmt.Println("密码验证失败！")
	}
}
```

### ✨ 核心特性

| 特性                | 描述                                                      |
|---------------------|-----------------------------------------------------------|
| **bcrypt 加密**     | 使用经过时间考验的 bcrypt 算法生成安全哈希                  |
| **自动盐值生成**    | 每次生成哈希时自动添加随机盐，提高密码安全性                 |
| **多样性哈希**      | 同一密码每次生成的哈希均不相同，防止哈希碰撞                |
| **简单易用**        | 提供直观的 API，方便在各类项目中集成和使用                    |
| **错误处理**        | 在哈希生成过程中提供详细的错误反馈                          |

### ⚠️ 注意事项

1. 请确保在生产环境中对错误信息进行妥善处理，避免泄露敏感信息。
2. bcrypt 默认成本为 8，若需要更高安全性可根据需要调整（注意更高成本会带来更高的计算消耗）。

### 🤝 参与贡献

- [贡献指南](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md)
- [提交 Issue](https://github.com/small-ek/antgo/issues)

[⬆ 返回顶部](#password---安全密码处理工具库--secure-password-utility)

---

## English

### 📖 Introduction

`password` is a straightforward and easy-to-use password utility library that provides functions for generating and verifying bcrypt-based password hashes.  
It is well-suited for user password storage, authentication systems, and any scenario that requires secure password handling.  
Each generated hash includes a random salt, version information, and configuration parameters, ensuring that even the same password produces a unique hash each time.

GitHub URL: [github.com/small-ek/antgo/utils/password](https://github.com/small-ek/antgo/utils/password)

### 📦 Installation

Install via `go get`:

```bash
go get github.com/small-ek/antgo/utils/password
```

### 🚀 Quick Start

#### Generating a Password Hash

Use the `Generate` function to create a secure bcrypt hash from a raw password.  
Example:

```go
package main

import (
	"fmt"
	"log"

	"github.com/small-ek/antgo/utils/password"
)

func main() {
	rawPassword := "mySecureP@ssw0rd"
	hashedPassword, err := password.Generate(rawPassword)
	if err != nil {
		log.Fatalf("Failed to generate password hash: %v", err)
	}
	fmt.Println("Generated hash:", hashedPassword)
}
```

#### Verifying a Password

Use the `Verify` function to check if a provided password matches the stored hash.  
Example:

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/utils/password"
)

func main() {
	storedHash := "$2a$08$u5hX9iIwBzr7rN7w5EjZne.Q/XO6zC8Q.f56kBvEUp6kT8YHvXGa6" // Example hash
	inputPassword := "mySecureP@ssw0rd"

	if password.Verify(storedHash, inputPassword) {
		fmt.Println("Password verification successful!")
	} else {
		fmt.Println("Password verification failed!")
	}
}
```

### ✨ Key Features

| Feature               | Description                                                      |
|-----------------------|------------------------------------------------------------------|
| **bcrypt Encryption** | Generates secure hashes using the well-established bcrypt algorithm |
| **Auto Salt Generation** | Automatically adds a random salt each time a hash is generated    |
| **Unique Hashes**     | The same password produces different hashes on each generation    |
| **Simplicity**        | Provides an intuitive API for easy integration into projects       |
| **Error Handling**    | Detailed error feedback during hash generation                     |

### ⚠️ Important Notes

1. Ensure that errors are handled appropriately in production environments to avoid exposing sensitive details.
2. The default cost for bcrypt is set to 8. Adjust the cost parameter as needed for increased security (note that higher costs incur higher computational overhead).

### 🤝 Contributing

- [Contribution Guide](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md)
- [Open an Issue](https://github.com/small-ek/antgo/issues)

[⬆ Back to Top](#password---secure-password-utility)