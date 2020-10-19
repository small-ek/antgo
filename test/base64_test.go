package test

import (
	"github.com/small-ek/ginp/encoding/base64"
	"log"
	"testing"
)

func TestBase64(t *testing.T) {
	for i := 0; i < 100; i++ {
		var str1 = base64.Encode([]byte("test"))
		log.Println(str1)
		var str2, err = base64.Decode(str1)
		log.Println(err)
		log.Println(str2)
	}
}
