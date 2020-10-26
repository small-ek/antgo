package test

import (
	"github.com/small-ek/ginp/net/ghttp"
	"log"
	"testing"
)

func TestHttp(t *testing.T) {
	for i := 0; i < 10; i++ {
		var result, err = ghttp.Client().SetBody(map[string]interface{}{
			"test": "test",
		}).Post("https://www.baidu.com")
		log.Println(string(result))
		log.Println(err)
	}
}
