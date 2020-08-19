package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

//密钥格式：PKCS#1
//密钥位数2048
var PrivateKeyJwt = []byte(`
-----BEGIN PRIVATE KEY-----
MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAMMBf8hUaDk+oCsQ
ecIJDrRt+EktQaGE+11Un1YWgjelgJR68JcS+2HsSTjkd+aUYFuy7uo7Jc0ugQqN
xTwPzyJXfIX0J3niM0MM4SNLiqGD+UYGJ3bbBiw33NPT4CQ0/ATKk7Y4kdGXWl9z
OMqV6YqWau88SpFAUENHj1rJrFKrAgMBAAECgYBEFSTo61dMDSpcfq8T6JeitPZH
ji5o1wXvqtjKdKdYCEdhD58qD62GnblezJ1z+n+95DX3v1jOTxssdRzUgGx/Ys13
Ukso6hzNgMOw46zwd8qoWuFuytWWc953FGsr+atWFrvfU8aEjBjWzhlTtFVRPaiS
42zS7OZzYdmCw/PJaQJBAPu9R23SePtKIufwkg9ejWKyfIwJfLfdrfPL3RSTryGy
ph6CahWcr19/fMRTJaTYQtfxnXc1quFow3X+IBgfMr0CQQDGTmibUOJGyVrjSNMe
dTntWaBjMdFRwNdl7EzgUuLePUFs0gbZ6SW5dMsK/3JfzyxYF3XRki7Lupju3Ano
RGWHAkAN/lCJJ0j4Vv+nuvSzjAL5+If51NEs+1KfGbb5XNhAXEjlq0QwXVxWR6Ts
2N5f0nGsxU6GgOI103gCCBVKoflVAkBvBzVwSE/4XAJEIOD7O50MM9Ml1p2gjTzM
Nwovyph0340C9XCajvvtIuQPq0gJNoBYbgIsLRGARWAc1BvD7I9/AkEAgQEnpQEI
isXUlyKSsakm+M+hzkoJxlizUiM3tN9cIfsIBXdWv9LNGRp2gl8Sa69ri3EdqQXv
0PcStMOn2IX1kw==
-----END PRIVATE KEY-----
`)
var PublicKeyJwt = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDDAX/IVGg5PqArEHnCCQ60bfhJ
LUGhhPtdVJ9WFoI3pYCUevCXEvth7Ek45HfmlGBbsu7qOyXNLoEKjcU8D88iV3yF
9Cd54jNDDOEjS4qhg/lGBid22wYsN9zT0+AkNPwEypO2OJHRl1pfczjKlemKlmrv
PEqRQFBDR49ayaxSqwIDAQAB
-----END PUBLIC KEY-----
`)

//创建Jwttoken
func Create(manifest map[string]interface{}) (tokenStr string) {
	Key, _ := jwt.ParseRSAPrivateKeyFromPEM(PrivateKeyJwt)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat":      time.Now().Unix(),                      // Token颁发时间
		"nbf":      time.Now().Unix(),                      // Token生效时间
		"exp":      time.Now().Add(time.Hour * 168).Unix(), // Token过期时间，目前是24小时
		"manifest": manifest,                               // 主题
	})
	result, err := token.SignedString(Key)
	if err != nil {
		log.Println("error:创建Jwt失败" + err.Error())
	}
	return result
}

//解密
func Decode(tokenStr string) (manifest map[string]interface{}, err error) {
	result := map[string]interface{}{}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(PublicKeyJwt)
	// 基于公钥验证Token合法性
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// 基于JWT的第一部分中的alg字段值进行一次验证
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("验证Token的加密类型错误")
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
	return result, errors.New("Token无效或者无对应值")
}
