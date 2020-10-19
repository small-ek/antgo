package base64

import (
	"encoding/base64"
	"log"
)

func Encode(str []byte) string {
	encodeString := base64.StdEncoding.EncodeToString(str)
	return encodeString
}

func Decode(str string) string {
	decodeBytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Println(err.Error())
	}
	return string(decodeBytes)
}
