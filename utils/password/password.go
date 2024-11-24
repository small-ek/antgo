package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// Generate 生成新密码,返回哈希密码每次生成的哈希密码都不一样
func Generate(password string) (string, error) {
	// 哈希密码
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(hashed), err
}

// Verify 验证密码是否正确
func Verify(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == nil {
		return true
	}
	return false
}
