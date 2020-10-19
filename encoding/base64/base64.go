package base64

import (
	"encoding/base64"
)

func Encode(str []byte) string {
	encodeString := base64.StdEncoding.EncodeToString(str)
	return encodeString
}

func Decode(str string) (string, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(str)
	return string(decodeBytes), err
}
