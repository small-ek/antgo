package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

//New...
type New struct {
	PrivateKey []byte
	PublicKey  []byte
}

//Default Default key
func Default(publicKey, privateKey []byte) *New {
	return &New{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
}

//Encrypt RSA encryption
func (this *New) Encrypt(origData string) (string, error) {
	block, _ := pem.Decode(this.PublicKey)
	if block == nil {
		return "", errors.New("Public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	pub := pubInterface.(*rsa.PublicKey)
	RsaEncrypt, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(origData))
	encodeString := base64.StdEncoding.EncodeToString(RsaEncrypt)
	return encodeString, err
}

//Decrypt RSA decryption
func (this *New) Decrypt(ciphertext string) (string, error) {
	decodeBytes, _ := base64.StdEncoding.DecodeString(ciphertext)
	block, _ := pem.Decode(this.PrivateKey)
	if block == nil {
		return "", errors.New("Decryption failed")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	RsaDecrypt, err := rsa.DecryptPKCS1v15(rand.Reader, priv, decodeBytes)
	return string(RsaDecrypt), err
}
