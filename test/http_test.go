package test

import (
	"github.com/small-ek/ginp/net/ghttp"
	"log"
	"testing"
)

func TestHttp(t *testing.T) {
	/*var result, err = ghttp.Client().SetProxy("http://127.0.0.1:58591").Get("https://www.facebook.com")*/
	for i := 0; i < 10; i++ {
		var result, err = ghttp.Client().Get("https://stage-admin.shoplitlive.com/admin/left_menu")
		log.Println(string(result))
		log.Println(err)
	}
}
