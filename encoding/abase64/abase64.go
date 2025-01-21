package abase64

import (
	"encoding/base64"
)

// Encode 使用标准Base64编码将字节切片转换为字符串
// Encode converts a byte slice to a standard Base64 encoded string
func Encode(data []byte) string {
	// 预计算编码后长度以减少内存分配 Pre-calculate encoded length to reduce memory allocation
	encodedLen := base64.StdEncoding.EncodedLen(len(data))

	// 直接编码到目标缓冲区 Encode directly to destination buffer
	dst := make([]byte, encodedLen)
	base64.StdEncoding.Encode(dst, data)

	return string(dst)
}

// Decode 将Base64编码字符串解码为原始字节切片
// Decodes a Base64 encoded string to the original byte slice
func Decode(encodedStr string) ([]byte, error) {
	// 将输入字符串转换为字节切片 Convert input string to byte slice
	src := []byte(encodedStr)

	// 预计算最大可能解码长度 Pre-calculate maximum possible decoded length
	maxDecodedLen := base64.StdEncoding.DecodedLen(len(src))

	// 直接解码到目标缓冲区，避免额外内存分配
	dst := make([]byte, maxDecodedLen)
	n, err := base64.StdEncoding.Decode(dst, src)
	if err != nil {
		return nil, err
	}

	// 返回实际解码数据部分 Avoid creating unnecessary slice
	return dst[:n], nil
}
