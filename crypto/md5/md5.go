package md5

import (
	"crypto/md5"
	"encoding/hex"
)

//MD5 encryption
func Create(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
