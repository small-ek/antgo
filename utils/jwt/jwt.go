package jwt

import (
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"sync"
	"time"
)

//JwtManager New JwtManager parameter
type JwtManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	mutex      sync.RWMutex
}

const defaultExp = time.Hour * 144

var jwtManager *JwtManager
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
		jwtManager = j
	})

	return jwtManager, nil

}

// SetPublicKey Set public key<设置公钥Key>
func (j *JwtManager) SetPublicKey(publicKey []byte) *JwtManager {
	if publicKey != nil && len(publicKey) > 0 {
		j.publicKey, _ = jwt.ParseRSAPublicKeyFromPEM(publicKey)
	}
	return j
}

// SetPrivateKey Set private <设置私钥Key>
func (j *JwtManager) SetPrivateKey(privateKey []byte) *JwtManager {
	if privateKey != nil && len(privateKey) > 0 {
		j.privateKey, _ = jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	}
	return j
}

// Encrypt json web token encryption<json web token 加密>
func (j *JwtManager) Encrypt(row map[string]interface{}, exp ...int64) (string, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	MapClaims := jwt.MapClaims{}
	nowTime := time.Now().Unix()
	//设置过期时间
	if len(exp) > 0 {
		row["exp"] = exp[0]
	} else {
		row["exp"] = time.Now().Add(defaultExp).Unix()
	}
	row["iat"] = nowTime //签发时间
	row["nbf"] = nowTime //生效时间
	MapClaims = row
	return jwt.NewWithClaims(jwt.SigningMethodRS256, MapClaims).SignedString(j.privateKey)
}

// Decode json web token decryption<json web token解密>
func (j *JwtManager) Decode(tokenStr string) (map[string]interface{}, error) {
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
