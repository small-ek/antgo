package test

import (
	"github.com/small-ek/ginp/net/ghttp"
	"log"
	"testing"
)

func TestHttp(t *testing.T) {
	var http = ghttp.Client().SetTimeout(10).SetProxy("http://127.0.0.1:58591")
	var result, err = http.Get("https://www.facebook.com/")
	for i, i2 := range http.GetHeader() {
		log.Println(i)
		log.Println(i2)
	}
	log.Println(string(result))
	log.Println(err)
}
