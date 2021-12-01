package test

import (
	"errors"
	"flag"
	"github.com/small-ek/antgo/conv"
	"github.com/small-ek/antgo/os/logger"
	"go.uber.org/zap"
	"log"
	"testing"
)

type Str struct {
	Name  string
	Name2 int
}

func TestLogger(t *testing.T) {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	flag.Parse()
	logger.Default("./log/ek2.log").SetServiceName("api").Register()
	logger.Info("123", `{"test":"1221"}`)
	logger.Write.Info("22333")
	data := conv.Bytes(Str{
		Name:  "123",
		Name2: 12,
	})
	logger.Write.Info("123", zap.ByteString("leve", data))

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
