package test

import (
	"encoding/hex"
	"github.com/small-ek/ginp/crypto/aes"
	"github.com/small-ek/ginp/encoding/base64"
	"log"
	"testing"
)

func TestCBC(t *testing.T) {
	origData := []byte("Hello World") // 待加密的数据
	key := []byte("ABCDEFGHIJKLMNOP") // 加密的密钥
	log.Println("原文：", string(origData))

	encrypted := aes.EncryptCBC(origData, key)
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", base64.Encode(encrypted))
	decrypted := aes.DecryptCBC(encrypted, key)
	log.Println("解密结果：", string(decrypted))
}

func TestCFB(t *testing.T) {
	origData := []byte("Hello World") // 待加密的数据
	key := []byte("ABCDEFGHIJKLMNOP") // 加密的密钥
	log.Println("原文：", string(origData))

	encrypted := aes.EncryptCFB(origData, key)
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", base64.Encode(encrypted))
	decrypted := aes.DecryptCFB(encrypted, key)
	log.Println("解密结果：", string(decrypted))
}
