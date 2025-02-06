# aemail - 轻量级邮件发送库 / Lightweight Email Sending Library

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`aemail` 是一个简洁高效的 Go 语言邮件发送库，支持 SMTP 发送邮件，提供 TLS 安全连接，并支持附件发送。它提供了一种简便的方式来在 Go 应用程序中快速集成邮件功能。

GitHub 地址: [github.com/small-ek/antgo/aemail](https://github.com/small-ek/antgo/aemail)

### 📦 安装

```bash
go get github.com/small-ek/antgo/aemail
```

### 🚀 快速开始

#### 创建邮件客户端

```go
package main

import (
	"github.com/small-ek/antgo/aemail"
)

func main() {
	mailer := aemail.NewMailer("your-email@example.com", "your-password")
}
```

#### 发送邮件

```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/aemail"
)

func main() {
	mailer := aemail.NewMailer("your-email@example.com", "your-password")
	email := &aemail.Email{
		To:      []string{"recipient@example.com"},
		Subject: "测试邮件",
		Text:    "这是邮件正文",
	}

	if err := mailer.Send(email); err != nil {
		fmt.Println("发送失败:", err)
	} else {
		fmt.Println("邮件发送成功")
	}
}
```

#### 发送 HTML 邮件

```go
email := &aemail.Email{
	To:      []string{"recipient@example.com"},
	Subject: "HTML 邮件",
	HTML:    "<h1>欢迎</h1><p>这是一封 HTML 格式的邮件</p>",
}
mailer.Send(email)
```

#### 发送带附件的邮件

```go
email := &aemail.Email{
	To:    []string{"recipient@example.com"},
	Subject: "附件测试",
	Text:    "请查看附件",
	Files:   []string{"./example.pdf"},
}
mailer.Send(email)
```

### 🔧 高级用法

#### 自定义 SMTP 配置

```go
mailer := aemail.NewMailer("your-email@example.com", "your-password").WithCustomSMTP("smtp.example.com", 587, true)
```

#### 快速发送文本邮件

```go
mailer.QuickSend([]string{"recipient@example.com"}, "快速邮件", "这是一封快速邮件")
```

#### 使用选项模式创建邮件

```go
email := aemail.NewEmail(
	aemail.WithTo("recipient@example.com"),
	aemail.WithSubject("选项模式邮件"),
	aemail.WithText("使用选项模式创建的邮件"),
)
mailer.Send(email)
```

### ✨ 核心特性

| 特性           | 描述                          |
| ------------ | --------------------------- |
| **TLS 支持**   | 通过 SMTP 进行安全邮件传输            |
| **附件支持**     | 允许发送带有多个附件的邮件               |
| **自定义 SMTP** | 可自由配置 SMTP 服务器地址、端口及 TLS 设置 |
| **快速发送**     | 通过 `QuickSend` 发送简短文本邮件     |
| **链式调用**     | 通过 `WithCustomSMTP` 进行自定义配置 |

### ⚠️ 注意事项

1. 确保 SMTP 服务器支持的端口和 TLS 配置正确。
2. 使用 QQ 邮箱时，需启用 SMTP 并使用授权码。
3. 发送失败时，请检查 SMTP 服务器地址、端口、TLS 设置。

### 🤝 参与贡献

[贡献指南](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [提交 Issue](https://github.com/small-ek/antgo/issues)

---

## English

### 📖 Introduction

`aemail` is a lightweight and efficient Go email sending library that supports SMTP with TLS security and attachment handling. It provides an easy way to integrate email functionality into Go applications.

GitHub URL: [github.com/small-ek/antgo/aemail](https://github.com/small-ek/antgo/aemail)

### 📦 Installation

```bash
go get github.com/small-ek/antgo/aemail
```

### 🚀 Quick Start

#### Create Mail Client

```go
mailer := aemail.NewMailer("your-email@example.com", "your-password")
```

#### Send Email

```go
email := &aemail.Email{
	To:      []string{"recipient@example.com"},
	Subject: "Test Email",
	Text:    "This is the email body",
}
mailer.Send(email)
```

#### Send HTML Email

```go
email := &aemail.Email{
	To:      []string{"recipient@example.com"},
	Subject: "HTML Email",
	HTML:    "<h1>Welcome</h1><p>This is an HTML email.</p>",
}
mailer.Send(email)
```

#### Send Email with Attachment

```go
email := &aemail.Email{
	To:    []string{"recipient@example.com"},
	Subject: "Attachment Test",
	Text:    "Please check the attachment",
	Files:   []string{"./example.pdf"},
}
mailer.Send(email)
```

### 🔧 Advanced Usage

#### Custom SMTP Configuration

```go
mailer := aemail.NewMailer("your-email@example.com", "your-password").WithCustomSMTP("smtp.example.com", 587, true)
```

#### Quick Send Text Email

```go
mailer.QuickSend([]string{"recipient@example.com"}, "Quick Email", "This is a quick email.")
```

#### Use Option Pattern for Email Creation

```go
email := aemail.NewEmail(
	aemail.WithTo("recipient@example.com"),
	aemail.WithSubject("Option Pattern Email"),
	aemail.WithText("Email created using option pattern"),
)
mailer.Send(email)
```

### ✨ Key Features

| Feature         | Description                           |
| --------------- | ------------------------------------- |
| **TLS Support** | Secure email transmission over SMTP   |
| **Attachments** | Supports sending multiple attachments |
| **Custom SMTP** | Configurable SMTP server, port, TLS   |
| **Quick Send**  | `QuickSend` for short text emails     |
| **Fluent API**  | `WithCustomSMTP` for custom config    |

### ⚠️ Important Notes

1. Ensure the SMTP server supports the correct port and TLS configuration.
2. When using QQ Mail, enable SMTP and use an authorization code.
3. If sending fails, check SMTP server settings and credentials.

### 🤝 Contributing
[Contribution Guide](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/small-ek/antgo/issues)

[⬆ Back to Top](#中文)