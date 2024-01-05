package main

import (
	"github.com/small-ek/antgo/net/ahttp"
	"github.com/small-ek/antgo/os/alog"
	"go.uber.org/zap"
	"log"
)

func main() {
	//log.SetFlags(log.Llongfile | log.LstdFlags)
	//flag.Parse()
	alog.Default("./log/ek2.log").SetServiceName("api").SetConsole(true).Register()
	//alog.Write.Info("1221", zap.String("12", "1221"))
	alog.Info("1111", zap.String("1212", "2222"), zap.String("1111", "2233"))
	var http = ahttp.Client()

	var result, err = http.Debug(false).Get("https://www.baidu.com/")
	log.Println(string(result))
	log.Println(err)
	var _, err2 = http.Debug(false).Get("https://weibo.com/newlogin")
	//log.Println(result2)
	log.Println(err2)
}
