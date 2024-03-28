package jwt

import (
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"sync"
	"time"
)

// New JwtManager parameter
type JwtManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	Exp        int64 //Expiration timestamp Default 15 days
	mutex      sync.Mutex
}

const defaultExp = time.Hour * 168

var jwtStr *JwtManager
var once sync.Once

// New function
func New(publicKey, privateKey []byte) (j *JwtManager, err error) {
	once.Do(func() {
		j = &JwtManager{}

		if privateKey != nil && len(privateKey) > 0 {
			j.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKey)
			if err != nil {
				return
			}
		}

		if publicKey != nil && len(publicKey) > 0 {
			j.publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKey)
			if err != nil {
				return
			}
		}
		jwtStr = j
	})

	return jwtStr, nil

}

// Encrypt json web token encryption<json web token 加密>
func (j *JwtManager) Encrypt(row map[string]interface{}, exp ...int64) (string, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	MapClaims := jwt.MapClaims{}
	if len(exp) > 0 {
		j.Exp = exp[0]
	} else {
		j.Exp = time.Now().Add(defaultExp).Unix()
	}
	MapClaims = row
	return jwt.NewWithClaims(jwt.SigningMethodRS256, MapClaims).SignedString(j.privateKey)
}

// Decode json web token decryption<json web token解密>
func (j *JwtManager) Decode(tokenStr string) (map[string]interface{}, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	result := map[string]interface{}{}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("token encryption type error")
		}
		return j.publicKey, nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result = claims
		return result, nil
	}
	return result, errors.New("token invalid")
}
