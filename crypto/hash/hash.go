package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash/crc32"
)

// Sha256 encryption
func Sha256(str string) string {
	var h = sha256.New()
	h.Write([]byte(str))
	sum := h.Sum(nil)
	var result = hex.EncodeToString(sum)
	return result
}

// Md5 encryption
func Md5(str string) string {
	var h = md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// Sha1 encryption
func Sha1(str string) string {
	var h = sha1.New()
	h.Write([]byte(str))
	var sum = h.Sum(nil)
	var result = hex.EncodeToString(sum)
	return result
}

// Crc32 crc32 encryption
func Crc32(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}
