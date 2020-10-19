package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type New struct {
	PrivateKey []byte
	PublicKey  []byte
	Exp        int64
}

func Default(PublicKey, PrivateKey []byte, exp ...int64) *New {
	var Exp = time.Now().Add(time.Hour * 360).Unix()
	if len(exp) > 0 {
		Exp = exp[0]
	}

	return &New{
		PublicKey:  PublicKey,
		PrivateKey: PrivateKey,
		Exp:        Exp,
	}
}

func (this *New) Encrypt(manifest map[string]interface{}) (string, error) {
	Key, _ := jwt.ParseRSAPrivateKeyFromPEM(this.PrivateKey)
	if this.Exp == 0 {
		this.Exp = time.Now().Add(time.Hour * 168).Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat":      time.Now().Unix(),
		"nbf":      time.Now().Unix(),
		"exp":      this.Exp,
		"manifest": manifest,
	})
	return token.SignedString(Key)
}

func (this *New) Decode(token_str string) (manifest map[string]interface{}, err error) {
	result := map[string]interface{}{}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(this.PublicKey)

	token, err := jwt.Parse(token_str, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("Token Encryption Type Error")
		}
		return publicKey, nil
	})

	if err != nil {
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result = claims["manifest"].(map[string]interface{})
		return result, nil
	}
	return result, errors.New("Token invalid")
}
