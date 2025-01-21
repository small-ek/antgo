package arsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

// Rsa 结构体包含公私钥信息
type Rsa struct {
	PublicKey  []byte
	PrivateKey []byte
}

// New 创建RSA实例
func New(publicKey, privateKey []byte) *Rsa {
	return &Rsa{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
}

// Encrypt RSA加密
func (r *Rsa) Encrypt(origData string) (string, error) {
	pub, err := parsePublicKey(r.PublicKey)
	if err != nil {
		return "", err
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(origData))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt RSA解密
func (r *Rsa) Decrypt(ciphertext string) (string, error) {
	priv, err := parsePrivateKey(r.PrivateKey)
	if err != nil {
		return "", err
	}

	decoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, priv, decoded)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// parsePublicKey 解析多种格式的公钥
func parsePublicKey(publicKey []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(publicKey)
	if block != nil {
		switch block.Type {
		case "RSA PUBLIC KEY": // PKCS#1
			return x509.ParsePKCS1PublicKey(block.Bytes)
		case "PUBLIC KEY": // PKCS#8
			pub, err := x509.ParsePKIXPublicKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			rsaPub, ok := pub.(*rsa.PublicKey)
			if !ok {
				return nil, errors.New("not an RSA public key")
			}
			return rsaPub, nil
		default:
			return nil, fmt.Errorf("unsupported PEM type: %s", block.Type)
		}
	}

	// 尝试DER格式解析
	if pub, err := x509.ParsePKCS1PublicKey(publicKey); err == nil {
		return pub, nil
	}

	pub, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("DER data does not contain RSA public key")
	}
	return rsaPub, nil
}

// parsePrivateKey 解析多种格式的私钥
func parsePrivateKey(privateKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKey)
	if block != nil {
		switch block.Type {
		case "RSA PRIVATE KEY": // PKCS#1
			return x509.ParsePKCS1PrivateKey(block.Bytes)
		case "PRIVATE KEY": // PKCS#8
			priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			rsaPriv, ok := priv.(*rsa.PrivateKey)
			if !ok {
				return nil, errors.New("not an RSA private key")
			}
			return rsaPriv, nil
		default:
			return nil, fmt.Errorf("unsupported PEM type: %s", block.Type)
		}
	}

	// 尝试DER格式解析
	if priv, err := x509.ParsePKCS1PrivateKey(privateKey); err == nil {
		return priv, nil
	}

	priv, err := x509.ParsePKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	rsaPriv, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("DER data does not contain RSA private key")
	}
	return rsaPriv, nil
}
