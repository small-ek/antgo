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
	var result, err = http.Debug().Get("https://www.baidu.com/")
	log.Println(result)
	log.Println(err)
	var result2, err2 = http.Debug().Get("https://www.baidu2.com/")
	log.Println(result2)
	log.Println(err2)
	for i := 0; i < 2; i++ {

	}
}
