package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

//Sha256 encryption
func Sha256(str string) string {
	var h = sha256.New()
	h.Write([]byte(str))
	sum := h.Sum(nil)
	var result = hex.EncodeToString(sum)
	return result
}

//MD5 encryption
func MD5(str string) string {
	var h = md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//Sha1 encryption
func Sha(str string) string {
	var h = sha1.New()
	h.Write([]byte(str))
	var sum = h.Sum(nil)
	var result = hex.EncodeToString(sum)
	return result
}
