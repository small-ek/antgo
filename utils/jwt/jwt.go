package jwt

import (
	"crypto/rsa"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/small-ek/antgo/os/aerror"
	"sync"
	"time"
)

// New Jwt parameter
type Jwt struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	Exp        int64 //Expiration timestamp Default 15 days
	mutex      sync.Mutex
}

const defaultExp = time.Hour * 168

// New function
func New(publicKey, privateKey []byte, exp ...int64) (*Jwt, error) {
	var err error
	j := &Jwt{}

	if len(exp) > 0 {
		j.Exp = exp[0]
	} else {
		j.Exp = time.Now().Add(defaultExp).Unix()
	}

	j.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return nil, aerror.WithMessage(err, "Private key error")
	}

	j.publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, aerror.WithMessage(err, "Public key error")
	}

	return j, nil

}

// Encrypt json web token encryption<json web token 加密>
func (j *Jwt) Encrypt(row map[string]interface{}) (string, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	MapClaims := jwt.MapClaims{}
	MapClaims = row
	return jwt.NewWithClaims(jwt.SigningMethodRS256, MapClaims).SignedString(j.privateKey)
}

// Decode json web token decryption<json web token解密>
func (j *Jwt) Decode(tokenStr string) (map[string]interface{}, error) {
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
