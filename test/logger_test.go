package test

import (
	"errors"
	"flag"
	"github.com/small-ek/ginp/os/logger"
	"log"
	"testing"
)

func TestLogger(t *testing.T) {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	flag.Parse()
	logger.Default("./log/ek2.log").SetServiceName("api").Register()
	logger.Info("123", `{"test":"1221"}`)
	logger.Write.Info("22333")
	test()
}
func test() {
	test23()
}
func test23() {
	var err = errors.New("")
	logger.Error(err)
	var err2 = errors.New("错误测试2")
	logger.Error(err2)
}
