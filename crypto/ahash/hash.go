package ahash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash/crc32"
)

// MD5 计算字符串的 MD5 哈希值
// MD5 calculates the MD5 hash value of a string
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	sum := h.Sum(nil)
	var buf [md5.Size * 2]byte // 32字节
	hex.Encode(buf[:], sum)
	return string(buf[:])
}

// SHA1 计算字符串的 SHA1 哈希值
// SHA1 calculates the SHA1 hash value of a string
func SHA1(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	sum := h.Sum(nil)
	var buf [sha1.Size * 2]byte // 40字节
	hex.Encode(buf[:], sum)
	return string(buf[:])
}

// SHA256 计算字符串的 SHA256 哈希值
// SHA256 calculates the SHA256 hash value of a string
func SHA256(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	sum := h.Sum(nil)
	var buf [sha256.Size * 2]byte // 64字节
	hex.Encode(buf[:], sum)
	return string(buf[:])
}

// SHA512 计算字符串的 SHA512 哈希值
// SHA512 calculates the SHA512 hash value of a string
func SHA512(str string) string {
	h := sha512.New()
	h.Write([]byte(str))
	sum := h.Sum(nil)
	var buf [sha512.Size * 2]byte // 128字节
	hex.Encode(buf[:], sum)
	return string(buf[:])
}

// Crc32 计算字符串的 CRC32 校验和
// Crc32 calculates the CRC32 checksum of a string
func Crc32(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}
