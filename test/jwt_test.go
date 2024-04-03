package test

import (
	"flag"
	"github.com/small-ek/antgo/os/alog"
	"github.com/small-ek/antgo/utils/jwt"
	"go.uber.org/zap"
	"log"
	"sync"
	"testing"
	"time"
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
	//var j, err = jwt.New(PublicKey, PrivateKey)
	//log.Println(err)
	//j = j.SetPublicKey(publicKey).SetPrivateKey(privateKey)
	//var token, err2 = j.Encrypt(data, time.Now().Add(time.Minute*1).Unix())
	//log.Println(token)
	//log.Println(err2)
	//var j2, err3 = jwt.New(PublicKey, nil)
	//log.Println(err3)
	//var getData, err4 = j2.Decode(token)
	//log.Println(getData)
	//log.Println(err4)
	//var getData2, err5 = j2.Decode("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTIyMTc3MDQsImlhdCI6MTcxMTY5OTMwNCwibmJmIjoxNzExNjk5MzA0LCJ0ZXN0IjoidGVzdCJ9.nzs_r92fY-Z24GX6LtMeQBdQm2B63tY4gbzufx0Z61nXPoLLXcA9dH5zmQpfQ00ivfJd5SxzxDwF_tHzYagZCcGuSsDnnXNNkyCKfF6e2A4s5jYbQAt39x4frimRMclmckq7ko1uCEkeRiNCsctYm5XHEOT_PTKvecHkiGFXnb8")
	//log.Println(getData2)
	//log.Println(err5)
	log.SetFlags(log.Llongfile | log.LstdFlags)
	flag.Parse()
	alog.New("./log/ek2.log").SetServiceName("api").Register()

	//if err != nil {
	//	t.Fatalf("Failed to create JwtManagerFactory: %v", err)
	//}

	var wg sync.WaitGroup
	numWorkers := 10000
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			jwtManagerFactory := jwt.New().SetPublicKey(PublicKey).SetPrivateKey(PrivateKey)
			jwtManager, err11 := jwtManagerFactory.Encrypt(data, time.Now().Add(time.Minute*1).Unix())

			jwtManager2, err12 := jwtManagerFactory.Decode(jwtManager)

			alog.Write.Info("123", zap.Any("jwtManager2", jwtManager2), zap.Error(err12), zap.Error(err11))
			// 在这里执行对 JwtManager 实例的操作
			// 例如，调用 Encrypt 或 Decode 方法
		}()
	}

	wg.Wait()
}
