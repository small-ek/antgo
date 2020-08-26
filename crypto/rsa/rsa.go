package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/pkg/errors"
)

//密钥格式：PKCS#1
//密钥位数1024
var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICWgIBAAKBgQChfYIOm7bOOlHF2VcEwqGaR1OxVT8fpJIjmvz6+k2EbiMo5e6b
L1pyRUXEVyXf+x06bJtXC1QCOKtwmnUKPAsSrWSP5jYmyLCjNaTJqwcCgyAXMA7k
AVoLRfjTxt4XYSm+P4HS+9w1yLZvsZ1XbRIhxgVHDtknkgANt7OREpBYowIDAQAB
An8BHMtcai9AZCJ9cQDcDnswQI0KnNAwxeCXIqrRD8yUsNBFwEnLtAnKucdBQqHb
cxIC9MaAZtoGmwIMZEl2BxmRYMhYZZyHdrKBNZJJgmV8Cz3Q1sRC12KJ1LMlSziV
8NO5tbKUOKzkY9GFItjgriDDfd1SAjNE+B7ydbN5cFKBAkEA/HTW3Xqt7Y8qxM51
4Za0RF1NQmVuyspKsaVcRCE/Igeq+T/pZ0jWvo/BFGDrFME4BTFbLJsP5URjZatH
wVzkIwJBAKPBy/xxzrpCwRVTMP2PkouqtUFYi8vYWK1j0DVPvNNkikYugOC+h6OQ
lbPiFspGgA1O2/rRqGzDT0fGdAjhwYECQQCwjeHKiMZkchB+DMmiF6xAd2PVwGxI
REsSi8vIFdw6J1Sp9cl8oxMTuCNW5iThofNUplzWCCeItlgxPST0lMszAkAhY9un
DsGbQw9BvOPJX+P+rIEm4NooZ2W1fRuwMyEKbX6wTr0illbr6AhOVHRXLEbh78l0
/Bj+jFh3ByUTxoyBAkAGZjTBaWJuQuNnD3Do+CuHAoJ3k9drAWEEAepQdzIIom0m
7+Vc8wQsdZSLK+ATqIM/KkC00dR1462axPXUR6f3
-----END RSA PRIVATE KEY-----
`)
var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQChfYIOm7bOOlHF2VcEwqGaR1Ox
VT8fpJIjmvz6+k2EbiMo5e6bL1pyRUXEVyXf+x06bJtXC1QCOKtwmnUKPAsSrWSP
5jYmyLCjNaTJqwcCgyAXMA7kAVoLRfjTxt4XYSm+P4HS+9w1yLZvsZ1XbRIhxgVH
DtknkgANt7OREpBYowIDAQAB
-----END PUBLIC KEY-----
`)

// RSA加密
func RsaEncrypt(origData string) (string, error) {
	block, _ := pem.Decode(publicKey) //将密钥解析成公钥实例
	if block == nil {
		return "", errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes) //解析pem.Decode（）返回的Block指针实例
	if err != nil {
		return "", err
	}
	pub := pubInterface.(*rsa.PublicKey)

	RsaEncrypt, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(origData)) //RSA算法加密
	//base64加密
	encodeString := base64.StdEncoding.EncodeToString(RsaEncrypt)

	return encodeString, err
}

// RSA解密
func RsaDecrypt(ciphertext string) (string, error) {
	//base64解密
	decodeBytes, _ := base64.StdEncoding.DecodeString(ciphertext)

	block, _ := pem.Decode(privateKey) //将密钥解析成私钥实例

	if block == nil {
		return "", errors.New("解密失败")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes) //解析pem.Decode（）返回的Block指针实例

	if err != nil {
		return "", err
	}

	RsaDecrypt, err := rsa.DecryptPKCS1v15(rand.Reader, priv, decodeBytes)

	return string(RsaDecrypt), err //RSA算法解密
}
