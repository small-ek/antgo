package test

import (
	"flag"
	"github.com/small-ek/antgo/os/logger"
	"github.com/small-ek/antgo/utils/conv"
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
}
