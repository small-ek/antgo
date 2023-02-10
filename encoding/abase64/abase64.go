package abase64

import (
	"encoding/base64"
)

// Encode base64 encoding
func Encode(str []byte) string {
	return base64.StdEncoding.EncodeToString(str)
}

// Decode base64 decoding
func Decode(str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(str)
}
