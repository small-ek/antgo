package test

import (
	"github.com/small-ek/antgo/encoding/abase64"
	"log"
	"testing"
)

func TestBase64(t *testing.T) {
	for i := 0; i < 100; i++ {
		var str1 = abase64.Encode([]byte("test"))
		log.Println(str1)
		var str2, _ = abase64.Decode(str1)
		log.Println(string(str2))
	}
}
