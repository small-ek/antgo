package test

import (
	"encoding/hex"
	"github.com/small-ek/antgo/crypto/aes"
	"github.com/small-ek/antgo/encoding/abase64"
	"log"
	"testing"
)

func TestCBC(t *testing.T) {
	origData := []byte("Hello World")
	key := []byte("ABCDEFGHIJKLMNOP")
	log.Println("原文：", string(origData))

	encrypted := aes.EncryptCBC(origData, key)
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", abase64.Encode(encrypted))
	decrypted := aes.DecryptCBC(encrypted, key)
	log.Println("解密结果：", string(decrypted))
}

func TestCFB(t *testing.T) {
	origData := []byte("Hello World")
	key := []byte("ABCDEFGHIJKLMNOP")
	log.Println("原文：", string(origData))

	encrypted := aes.EncryptCFB(origData, key)
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", abase64.Encode(encrypted))
	decrypted := aes.DecryptCFB(encrypted, key)
	log.Println("解密结果：", string(decrypted))
}
