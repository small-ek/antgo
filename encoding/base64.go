package encoding

import (
	"encoding/base64"
	"log"
)

// 使用BASE64算法对字节进行编码(Encode bytes using BASE64 algorithm)
// @str 编码参数(Encoding parameters)
func Encode(str [] byte) string {
	encodeString := base64.StdEncoding.EncodeToString(str)
	return encodeString
}

// 使用BASE64算法解码字符串(Use BASE64 algorithm to decode string)
// @str 解码参数(Decoding parameters)
func Decode(str string) string {
	decodeBytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Print("err:base64 decoding failed" + err.Error())
	}
	return string(decodeBytes)
}
