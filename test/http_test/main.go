package main

import (
	"github.com/small-ek/antgo/net/ahttp"
	"github.com/small-ek/antgo/os/logger"
	"log"
)

func main() {
	logger.Default("E:/go/src/antgo/test/log/ek2.log").SetServiceName("api").Register()
	var http = ahttp.Client()

	var _, err = http.Debug(false).Get("https://www.baidu.com/")
	//log.Println(result)
	log.Println(err)
	var _, err2 = http.Debug(false).Get("https://weibo.com/newlogin")
	//log.Println(result2)
	log.Println(err2)
}
