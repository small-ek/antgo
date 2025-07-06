package ahash

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"hash/crc32"
)

// MD5 computes the MD5 hash and returns hex-encoded string
// MD5哈希计算，返回十六进制编码字符串
func MD5(data string) string {
	sum := md5.Sum([]byte(data))
	return hex.EncodeToString(sum[:])
}

// SHA1 computes the SHA-1 hash and returns hex-encoded string
// SHA1哈希计算，返回十六进制编码字符串
func SHA1(data string) string {
	sum := sha1.Sum([]byte(data))
	return hex.EncodeToString(sum[:])
}

// SHA256 computes the SHA-256 hash and returns hex-encoded string
// SHA256哈希计算，返回十六进制编码字符串
func SHA256(data string) string {
	sum := sha256.Sum256([]byte(data))
	return hex.EncodeToString(sum[:])
}

// SHA512 computes the SHA-512 hash and returns hex-encoded string
// SHA512哈希计算，返回十六进制编码字符串
func SHA512(data string) string {
	sum := sha512.Sum512([]byte(data))
	return hex.EncodeToString(sum[:])
}

// CRC32 computes IEEE CRC32 checksum
// 计算CRC32 IEEE校验值
func CRC32(data string) uint32 {
	return crc32.ChecksumIEEE([]byte(data))
}

// SignSHA1 generates HMAC-SHA1 signature with base64 encoding
//
// Args:
//
//	data: message to sign 待签名数据
//	key: secret key 密钥
//
// Returns:
//
//	Base64-encoded signature Base64编码的签名结果
func SignSHA1(data, key string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// SignSHA256 generates HMAC-SHA256 signature with base64 encoding
//
// Args:
//
//	data: message to sign 待签名数据
//	key: secret key 密钥
//
// Returns:
//
//	Base64-encoded signature Base64编码的签名结果
func SignSHA256(data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// SignSHA512 generates HMAC-SHA512 signature with base64 encoding
//
// Args:
//
//	data: message to sign 待签名数据
//	key: secret key 密钥
//
// Returns:
//
//	Base64-encoded signature Base64编码的签名结果
func SignSHA512(data, key string) string {
	mac := hmac.New(sha512.New, []byte(key))
	mac.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
