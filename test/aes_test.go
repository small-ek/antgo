package test

import (
	"encoding/hex"
	"github.com/small-ek/antgo/crypto/aaes"
	"github.com/small-ek/antgo/encoding/abase64"
	"log"
	"testing"
)

func TestCBC(t *testing.T) {
	origData := []byte("Hello World")
	key := []byte("ABCDEFGHIJKLMNOP")
	log.Println("原文：", string(origData))

	encrypted := aaes.EncryptCBC(origData, key)
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", abase64.Encode(encrypted))
	decrypted := aaes.DecryptCBC(encrypted, key)
	log.Println("解密结果：", string(decrypted))
}

func TestCFB(t *testing.T) {
	origData := []byte("Hello World")
	key := []byte("ABCDEFGHIJKLMNOP")
	log.Println("原文：", string(origData))

	encrypted := aaes.EncryptCFB(origData, key)
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", abase64.Encode(encrypted))
	decrypted := aaes.DecryptCFB(encrypted, key)
	log.Println("解密结果：", string(decrypted))
}
