# arsa - RSA 加密解密库 / RSA Encryption and Decryption Utilities

[中文](#中文) | [English](#english)

---

## 中文

### 📖 简介

`arsa` 是一个基于 Go 语言的 RSA 加密解密库，支持公钥加密和私钥解密操作。适用于数据加密传输、数字签名、安全通信等场景。

GitHub地址: [github.com/small-ek/antgo/crypto/arsa](https://github.com/small-ek/antgo/crypto/arsa)

### 📦 安装

```bash
go get github.com/small-ek/antgo/crypto/arsa
```

### 🚀 快速开始

#### 初始化 RSA 实例
```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/crypto/arsa"
)

func main() {
	// 公钥和私钥（PEM 格式）
	publicKey := []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...`)
	privateKey := []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA...`)

	// 创建 RSA 实例
	rsaInstance := arsa.New(publicKey, privateKey)
}
```

#### 加密数据
```go
func main() {
	// 初始化 RSA 实例
	rsaInstance := arsa.New(publicKey, privateKey)

	// 加密数据
	plaintext := "需要加密的数据"
	ciphertext, err := rsaInstance.Encrypt(plaintext)
	if err != nil {
		fmt.Println("加密失败:", err)
		return
	}
	fmt.Println("加密后的数据:", ciphertext)
}
```

#### 解密数据
```go
func main() {
	// 初始化 RSA 实例
	rsaInstance := arsa.New(publicKey, privateKey)

	// 解密数据
	decryptedText, err := rsaInstance.Decrypt(ciphertext)
	if err != nil {
		fmt.Println("解密失败:", err)
		return
	}
	fmt.Println("解密后的数据:", decryptedText)
}
```

### 🔧 高级用法

#### 支持多种密钥格式
`arsa` 支持解析多种格式的公钥和私钥，包括 PKCS#1 和 PKCS#8 格式的 PEM 文件以及 DER 格式的密钥。

```go
func main() {
	// PKCS#8 格式的公钥
	publicKey := []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...`)

	// PKCS#1 格式的私钥
	privateKey := []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA...`)

	rsaInstance := arsa.New(publicKey, privateKey)
}
```

### ✨ 核心特性

| 特性                | 描述                                                                 |
|---------------------|--------------------------------------------------------------------|
| **多格式支持**       | 支持 PKCS#1、PKCS#8 格式的 PEM 文件以及 DER 格式的密钥             |
| **简单易用**         | 提供简洁的 API，快速实现 RSA 加密解密                              |
| **高性能**           | 基于 Go 语言原生 RSA 库实现，性能优异                             |
| **跨平台**           | 支持所有 Go 语言支持的平台                                         |

### ⚠️ 注意事项
1. 请确保公钥和私钥匹配，否则解密会失败。
2. RSA 加密的数据长度受密钥长度限制，建议加密较短的数据或使用对称加密结合 RSA 加密的方式。
3. 密钥文件需妥善保管，避免泄露。

### 🤝 参与贡献
[贡献指南](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [提交Issue](https://github.com/small-ek/antgo/issues)

---

## English

### 📖 Introduction

`arsa` is a Go-based RSA encryption and decryption library that supports public key encryption and private key decryption. It is suitable for scenarios such as secure data transmission, digital signatures, and secure communication.

GitHub URL: [github.com/small-ek/antgo/crypto/arsa](https://github.com/small-ek/antgo/crypto/arsa)

### 📦 Installation

```bash
go get github.com/small-ek/antgo/crypto/arsa
```

### 🚀 Quick Start

#### Initialize RSA Instance
```go
package main

import (
	"fmt"
	"github.com/small-ek/antgo/crypto/arsa"
)

func main() {
	// Public and private keys (PEM format)
	publicKey := []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...`)
	privateKey := []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA...`)

	// Create RSA instance
	rsaInstance := arsa.New(publicKey, privateKey)
}
```

#### Encrypt Data
```go
func main() {
	// Initialize RSA instance
	rsaInstance := arsa.New(publicKey, privateKey)

	// Encrypt data
	plaintext := "data to encrypt"
	ciphertext, err := rsaInstance.Encrypt(plaintext)
	if err != nil {
		fmt.Println("Encryption failed:", err)
		return
	}
	fmt.Println("Encrypted data:", ciphertext)
}
```

#### Decrypt Data
```go
func main() {
	// Initialize RSA instance
	rsaInstance := arsa.New(publicKey, privateKey)

	// Decrypt data
	decryptedText, err := rsaInstance.Decrypt(ciphertext)
	if err != nil {
		fmt.Println("Decryption failed:", err)
		return
	}
	fmt.Println("Decrypted data:", decryptedText)
}
```

### 🔧 Advanced Usage

#### Support for Multiple Key Formats
`arsa` supports parsing multiple formats of public and private keys, including PKCS#1 and PKCS#8 PEM files, as well as DER-encoded keys.

```go
func main() {
	// PKCS#8 format public key
	publicKey := []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...`)

	// PKCS#1 format private key
	privateKey := []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA...`)

	rsaInstance := arsa.New(publicKey, privateKey)
}
```

### ✨ Key Features

| Feature             | Description                                                     |
|---------------------|-----------------------------------------------------------------|
| **Multi-format**    | Supports PKCS#1, PKCS#8 PEM files, and DER-encoded keys         |
| **Easy to Use**     | Provides a simple API for quick RSA encryption and decryption   |
| **High Performance**| Built on Go's native RSA libraries for excellent performance    |
| **Cross-platform**  | Supports all platforms compatible with Go                       |

### ⚠️ Important Notes
1. Ensure that the public and private keys match; otherwise, decryption will fail.
2. The length of data encrypted with RSA is limited by the key size. It is recommended to encrypt short data or use a combination of symmetric encryption and RSA.
3. Keep key files secure to avoid leakage.

### 🤝 Contributing
[Contribution Guide](https://github.com/small-ek/antgo/blob/main/CONTRIBUTING.md) | [Open an Issue](https://github.com/small-ek/antgo/issues)

[⬆ Back to Top](#中文)