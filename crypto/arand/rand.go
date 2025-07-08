package arand

import (
	"crypto/rand"
	"errors"
)

const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomString 生成指定长度的随机字节串，字符来自预定义字符集
// 使用 crypto/rand 保证密码学安全，单次读取优化性能
func RandomString(length int) ([]byte, error) {
	if length <= 0 {
		return nil, errors.New("length must be positive integer")
	}

	// 直接复用结果缓冲区，避免额外分配
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}

	// 就地修改缓冲区，避免内存拷贝
	setSize := byte(len(charSet))
	for i := range buf {
		// 取模映射到字符集索引
		buf[i] = charSet[buf[i]%setSize]
	}

	return buf, nil
}
