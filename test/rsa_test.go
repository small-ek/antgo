package test

import (
	"fmt"
	"github.com/small-ek/antgo/crypto/arsa"
	"log"
	"testing"
)

var privateKey = []byte(`
MIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBAI0zk42yxBk507t1/gjLmRg99lb0b2yRi+eA5ta2TAIUZ6K42Kaf5OArHn0hZFRZoZPehvecD+EEY+BxBi1AxRsDzrSj1vcQHLUpE/ZJTDY2mo85f+wT+AEvHBzFO+dnhzuhiqbtmdtyZrSXTa6PhMfAQq7Fho5dYXvlIvMuhbXXAgMBAAECgYB5981DHuHo8FyGlztiGYwatpstLWN6Ipb42A2N9lXdjgafNpBeDcvEHzrs14U0n1/EvPlUtUe6FPK5EqhFOqeGhoAAJ0pq3sVx3BLYUK0RYPKoLs5qUwdyo4k81sV4r5mr00RIl2OcYG5TstTskg5whQSNc0cJf474CXnL0LzxMQJBAPFX1a6IwaxR+9/b7HD438L5T/Kdpw4e0rnqiUkUiz9hHyyzuidE0Har+xj88o6/wdVUzpyjJsWs9kovTKy9fy0CQQCVxtaIOBp8kBrEoC6oW26gizy4QtAVAfkS2ulcZc33R9ux19bKF/6jOjQ194YTg5NfpDwIartJRUKhQ/VTccuTAkAR8uBXbKBKuoYq7eY1uKybiYMingrwh+ZQIVs4bii0+/ofjvZHOVzvlnbEMvuvFh/KR9Zd29xkUyq19bKUHju5AkB2h1zPgFa1rPUCFiHWakUqGAZ9a6JwfZc3LLbwwEA3KU7bdwwr8sE5O56F9tTMLJw8XCSGJLECUyVfqgBDgRKxAkByYC36fUpIVtTZxBq++yIvFkgjxLrTlwkCJti9vMB3LsJdQg+3Ms0AilxIyd+wMRqD/h6L7QnQxIg3kAyf2xxw
`)
var publicKey = []byte(`
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCNM5ONssQZOdO7df4Iy5kYPfZW9G9skYvngObWtkwCFGeiuNimn+TgKx59IWRUWaGT3ob3nA/hBGPgcQYtQMUbA860o9b3EBy1KRP2SUw2NpqPOX/sE/gBLxwcxTvnZ4c7oYqm7Znbcma0l02uj4THwEKuxYaOXWF75SLzLoW11wIDAQAB
`)

func TestRSA(t *testing.T) {
	/*var rsa=rsa.Default()*/
	var result, err = arsa.New(publicKey, privateKey)
	fmt.Println(err)
	var str1, _ = result.Encrypt("admin")
	log.Println(str1)
	log.Println(result.Decrypt(str1))
}
