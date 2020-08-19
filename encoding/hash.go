package encoding

import (
	"crypto/sha256"
	"encoding/hex"
)

//使用sha256哈希函数
func Sha256(str string) string {
	//使用sha256哈希函数
	h := sha256.New()
	h.Write([]byte(str))
	sum := h.Sum(nil)

	//由于是十六进制表示，因此需要转换
	s := hex.EncodeToString(sum)
	return s
}
