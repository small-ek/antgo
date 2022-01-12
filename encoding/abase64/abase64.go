package abase64

import (
	"encoding/base64"
)

//Encode base64 encoding
func Encode(str []byte) string {
	encodeString := base64.StdEncoding.EncodeToString(str)
	return encodeString
}

//Decode base64 decoding
func Decode(str string) string {
	decodeBytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		panic(err)
	}
	return string(decodeBytes)
}
