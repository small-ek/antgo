package test

import (
	"github.com/small-ek/antgo/encoding/aurl"
	"log"
	"testing"
)

var urlStr string = `https://golang.org/x/crypto?go-get=1 +`
var src string = `http://username:password@hostname:9090/path?arg=value#anchor`

func TestUrl(t *testing.T) {

	var result = aurl.Encode(urlStr)
	log.Println(result)
	var result2, _ = aurl.Decode(result)
	log.Println(result2)
	var component = 0
	var result3, _ = aurl.ParseURL(src, component)
	log.Println(result3)
}
