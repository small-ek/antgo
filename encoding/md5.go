package encoding

import (
	"crypto/md5"
	"encoding/hex"
)

//MD5加密
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
