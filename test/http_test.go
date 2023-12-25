package test

import (
	"github.com/small-ek/antgo/net/ahttp"
	"log"
	"testing"
)

func TestHttp(t *testing.T) {
	var http = ahttp.Client()
	//var result, err = http.Debug().SetFile("test.jpg", "file").SetBody(map[string]interface{}{"name": "123.jpg"}).PostForm("http://127.0.0.1:102/upload_file")
	//log.Println(http)
	//
	//log.Println(string(result))
	//log.Println(err)

	for i := 0; i < 1; i++ {
		var result, err = http.Debug().Get("https://www.baidu.com/")
		log.Println(result)
		log.Println(err)
	}
}
