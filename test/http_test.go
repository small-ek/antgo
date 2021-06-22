package test

import (
	"github.com/small-ek/antgo/net/ahttp"
	"log"
	"testing"
)

func TestHttp(t *testing.T) {
	var http = ahttp.Client()
	var result, err = http.SetFile("test.jpg").SetFileKey("file").SetFileName("img.jpg").SetBody(map[string]interface{}{"name": "123.jpg"}).PostForm("http://127.0.0.1:102/upload_file")
	log.Println(http.GetResponse().Request)
	//for i, i2 := range http.GetHeader() {
	//	log.Println(i)
	//	log.Println(i2)
	//}
	log.Println(string(result))
	log.Println(err)
	//ahttp.Test("http://127.0.0.1:102/upload_file")
}
