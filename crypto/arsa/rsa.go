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

type Rsa struct {
	pub  *rsa.PublicKey
	priv *rsa.PrivateKey
}

// New 创建一个 Rsa 实例，会自动解析并缓存公私钥，提高性能
func New(publicKeyPEM, privateKeyPEM []byte) (*Rsa, error) {
	r := &Rsa{}

	if len(publicKeyPEM) > 0 {
		pub, err := parsePublicKey(publicKeyPEM)
		if err != nil {
			return nil, fmt.Errorf("parse public key failed: %w", err)
		}
		r.pub = pub
	}

	if len(privateKeyPEM) > 0 {
		priv, err := parsePrivateKey(privateKeyPEM)
		if err != nil {
			return nil, fmt.Errorf("parse private key failed: %w", err)
		}
		r.priv = priv
	}

	return r, nil
}

// Encrypt 对明文进行 RSA 加密，返回 Base64 编码密文
func (r *Rsa) Encrypt(plaintext string) (string, error) {
	if r.pub == nil {
		return "", errors.New("public key is not initialized")
	}

	cipherBytes, err := rsa.EncryptPKCS1v15(rand.Reader, r.pub, []byte(plaintext))
	if err != nil {
		return "", fmt.Errorf("encrypt error: %w", err)
	}

	return base64.StdEncoding.EncodeToString(cipherBytes), nil
}

// Decrypt 对 Base64 编码的密文进行 RSA 解密，返回明文
func (r *Rsa) Decrypt(ciphertext string) (string, error) {
	if r.priv == nil {
		return "", errors.New("private key is not initialized")
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("base64 decode error: %w", err)
	}

	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, r.priv, data)
	if err != nil {
		return "", fmt.Errorf("decrypt error: %w", err)
	}

	return string(plaintext), nil
}

// parsePublicKey 支持 PEM / DER / Base64 格式的公钥
func parsePublicKey(data []byte) (*rsa.PublicKey, error) {
	// PEM 格式
	if block, _ := pem.Decode(data); block != nil {
		switch block.Type {
		case "PUBLIC KEY":
			key, err := x509.ParsePKIXPublicKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			if pub, ok := key.(*rsa.PublicKey); ok {
				return pub, nil
			}
			return nil, errors.New("not RSA public key")
		case "RSA PUBLIC KEY":
			return x509.ParsePKCS1PublicKey(block.Bytes)
		}
	}

	// Base64 尝试
	if decoded, err := base64.StdEncoding.DecodeString(string(data)); err == nil {
		if pub, err := x509.ParsePKIXPublicKey(decoded); err == nil {
			if rsaPub, ok := pub.(*rsa.PublicKey); ok {
				return rsaPub, nil
			}
		}
		if pub, err := x509.ParsePKCS1PublicKey(decoded); err == nil {
			return pub, nil
		}
	}

	// DER 原始尝试
	if pub, err := x509.ParsePKCS1PublicKey(data); err == nil {
		return pub, nil
	}
	if pub, err := x509.ParsePKIXPublicKey(data); err == nil {
		if rsaPub, ok := pub.(*rsa.PublicKey); ok {
			return rsaPub, nil
		}
	}

	return nil, errors.New("invalid RSA public key")
}

// parsePrivateKey 支持 PEM / DER / Base64 格式的私钥
func parsePrivateKey(data []byte) (*rsa.PrivateKey, error) {
	// PEM 格式
	if block, _ := pem.Decode(data); block != nil {
		switch block.Type {
		case "PRIVATE KEY":
			key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			if priv, ok := key.(*rsa.PrivateKey); ok {
				return priv, nil
			}
			return nil, errors.New("not RSA private key")
		case "RSA PRIVATE KEY":
			return x509.ParsePKCS1PrivateKey(block.Bytes)
		}
	}

	// Base64 尝试
	if decoded, err := base64.StdEncoding.DecodeString(string(data)); err == nil {
		if priv, err := x509.ParsePKCS1PrivateKey(decoded); err == nil {
			return priv, nil
		}
		if key, err := x509.ParsePKCS8PrivateKey(decoded); err == nil {
			if rsaPriv, ok := key.(*rsa.PrivateKey); ok {
				return rsaPriv, nil
			}
		}
	}

	// DER 原始尝试
	if priv, err := x509.ParsePKCS1PrivateKey(data); err == nil {
		return priv, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(data); err == nil {
		if rsaPriv, ok := key.(*rsa.PrivateKey); ok {
			return rsaPriv, nil
		}
	}

	return nil, errors.New("invalid RSA private key")
}
