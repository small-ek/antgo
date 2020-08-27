package test

import (
	"github.com/small-ek/ginp/crypto/des"
	"github.com/small-ek/ginp/encoding/base64"
	"log"
	"testing"
)

func TestEcb(t *testing.T)  {
	key := []byte("11111111")
	text := []byte("12345678")
	padding := 0
	/*result := "858b176da8b12503"*/
	for i:=0;i<100;i++{
		var str,_=des.EncryptECB(text, key, padding)
		log.Println(base64.Encode(str))
		var str2,_=des.DecryptECB(str,key,padding)
		log.Println(string(str2))

	}
}

func TestCbc(t *testing.T)  {
	key := []byte("11111111")
	text := []byte("1234567812345678")
	padding := 0
	iv := []byte("12345678")

	for i:=0;i<100;i++{
		var str,_=des.EncryptCBC(text, key, iv, padding)
		log.Println(base64.Encode(str))
		var str2,_=des.DecryptCBC(str,key, iv, padding)
		log.Println(string(str2))
	}
}