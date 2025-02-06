package password

import (
	"golang.org/x/crypto/bcrypt"
)

// defaultCost 表示默认的哈希计算成本
// defaultCost represents the default hashing cost
const defaultCost = 8

// Generate 根据输入密码生成安全的bcrypt哈希
// 特征：
//   - 自动生成随机盐值
//   - 包含版本信息和哈希配置参数
//   - 每次调用生成不同的哈希值
//
// Generate generates secure bcrypt hash from password
// Features:
//   - Auto-generates random salt
//   - Contains version information and hash parameters
//   - Produces different hash values on each invocation
//
// 参数:
//
//	password: 需要哈希的原始密码 (raw password to be hashed)
//
// 返回:
//
//	string: 生成的哈希字符串 (generated hash string)
//	error: 错误信息 (error information)
func Generate(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), defaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// Verify 验证密码是否匹配存储的哈希值
// Verify verifies if password matches stored hash
//
// 参数:
//
//	hashedPassword: 之前生成的哈希密码 (previously generated password hash)
//	password:       需要验证的原始密码 (raw password to verify)
//
// 返回:
//
//	bool: 验证结果 (verification result)
func Verify(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
	return err == nil
}
