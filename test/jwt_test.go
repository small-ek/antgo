package test

import (
	"github.com/small-ek/ginp/jwt"
	"log"
	"testing"
)

func TestJwt(t *testing.T) {
	var data = map[string]interface{}{
		"test": "test",
	}
	var jwt = jwt.Default()
	var result2 = jwt.Encrypt(data)
	log.Println(result2)
	var getData, _ = jwt.Decode(result2)
	log.Println(getData)
}
