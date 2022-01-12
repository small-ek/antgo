package test

import (
	"github.com/small-ek/antgo/jwt"
	"log"
	"testing"
)

func TestJwt(t *testing.T) {

	var PrivateKey = []byte(`
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
	var PublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDDAX/IVGg5PqArEHnCCQ60bfhJ
LUGhhPtdVJ9WFoI3pYCUevCXEvth7Ek45HfmlGBbsu7qOyXNLoEKjcU8D88iV3yF
9Cd54jNDDOEjS4qhg/lGBid22wYsN9zT0+AkNPwEypO2OJHRl1pfczjKlemKlmrv
PEqRQFBDR49ayaxSqwIDAQAB
-----END PUBLIC KEY-----
`)

	var data = map[string]interface{}{
		"test": "test",
	}
	var j = jwt.New(PublicKey, PrivateKey, 1642648770)
	var token, err = j.Encrypt(data)
	log.Println(token)
	log.Println(err)

	var getData, err2 = j.Decode(token)
	log.Println(getData)
	log.Println(err2)
}
