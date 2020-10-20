package base64

import (
	"encoding/base64"
)

//Encode base64 encoding
func Encode(str []byte) string {
	encodeString := base64.StdEncoding.EncodeToString(str)
	return encodeString
}

//Decode base64 decoding
func Decode(str string) (string, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(str)
	return string(decodeBytes), err
}
