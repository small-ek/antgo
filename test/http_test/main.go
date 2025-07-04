package main

import (
	"context"
	"github.com/small-ek/antgo/os/alog"
	"go.uber.org/zap"
)

func main() {
	//log.SetFlags(log.Llongfile | log.LstdFlags)
	//flag.Parse()
	alog.New("./log").SetServiceName("api").SetConsole(true).Register()
	//alog.Write.Info("1221", zap.String("12", "1221"))
	alog.Info(context.Background(), "1111", zap.String("1212", "2222"), zap.String("1111", "2233"))
	alog.Write.Info("1111", zap.String("1212", "2222"), zap.String("1111", "2233"))
}
