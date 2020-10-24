package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

//New Jwt parameter
type New struct {
	PrivateKey []byte //Private key
	PublicKey  []byte //Public key
	Exp        int64  //Expiration timestamp
}

//Default function
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

//Encrypt json web token encryption<json web token 加密>
func (get *New) Encrypt(manifest map[string]interface{}) (string, error) {
	Key, _ := jwt.ParseRSAPrivateKeyFromPEM(get.PrivateKey)
	if get.Exp == 0 {
		get.Exp = time.Now().Add(time.Hour * 168).Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat":      time.Now().Unix(),
		"nbf":      time.Now().Unix(),
		"exp":      get.Exp,
		"manifest": manifest,
	})
	return token.SignedString(Key)
}

//Decode json web token decryption<json web token解密>
func (get *New) Decode(tokenStr string) (manifest map[string]interface{}, err error) {
	result := map[string]interface{}{}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(get.PublicKey)

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("token encryption type error")
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
	return result, errors.New("token invalid")
}
