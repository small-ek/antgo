package test

import (
	"github.com/small-ek/ginp/crypto/aes"
	"github.com/small-ek/ginp/encoding/base64"
	"log"
	"testing"
)

func TestAES(t *testing.T) {
	var content = []byte("pibigstar")
	var key_16 = []byte("1234567891234567")
	for i := 0; i < 100; i++ {
		var data, _ = aes.Encrypt(content, key_16)
		log.Println(base64.Encode(data))
		var data2, _ = aes.Decrypt(data, key_16)
		log.Println(string(data2))
	}
}
