package arsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

//New ...
type Rsa struct {
	PrivateKey []byte
	PublicKey  []byte
}

//New Default key
func New(publicKey, privateKey []byte) *Rsa {
	return &Rsa{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
}

//Encrypt RSA encryption
func (r *Rsa) Encrypt(origData string) (string, error) {
	block, _ := pem.Decode(r.PublicKey)
	if block == nil {
		return "", errors.New("public key error")
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
func (r *Rsa) Decrypt(ciphertext string) (string, error) {
	decodeBytes, _ := base64.StdEncoding.DecodeString(ciphertext)
	block, _ := pem.Decode(r.PrivateKey)
	if block == nil {
		return "", errors.New("decryption failed")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	RsaDecrypt, err := rsa.DecryptPKCS1v15(rand.Reader, priv, decodeBytes)
	return string(RsaDecrypt), err
}
