package test

import (
	"github.com/small-ek/ginp/ghttp"
	"log"
	"testing"
)

func TestHttp(t *testing.T) {
	var result, err = ghttp.Client().SetBody(map[string]interface{}{
		"test": "test",
	}).Post("https://www.baidu.com")
	log.Println(string(result))
	log.Println(err)
}
