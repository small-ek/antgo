package test

import (
	"flag"
	"fmt"
	"github.com/small-ek/antgo/frame/ant"
	"github.com/small-ek/antgo/net/ahttp"
	"github.com/small-ek/antgo/os/alog"
	"log"
	"testing"
)

var http = ahttp.New(nil).SetLog(ant.Log()).Client()

func TestHttp(t *testing.T) {
	RegisterLog()
	result, err := http.EnableGenerateCurlOnDebug().Get("https://www.baidu.com")
	fmt.Println(string(result.Body()))
	fmt.Println(err)
}

func RegisterLog() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	flag.Parse()
	alog.New("./log/ek2.log").SetServiceName("api").Register()
}
