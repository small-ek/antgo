package test

import (
	"flag"
	"github.com/small-ek/antgo/os/alog"
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
	alog.New("./log/app.log").SetServiceName("prod").Register()

	for i := 0; i < 1000000; i++ {
		alog.Write.Info("22333")
	}
	data := conv.Bytes(Str{
		Name:  "123",
		Name2: 12,
	})
	alog.Write.Info("123", zap.ByteString("leve", data))
}
