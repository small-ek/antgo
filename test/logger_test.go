package test

import (
	"github.com/small-ek/ginp/os/logger"
	"testing"
)

func TestLogger(t *testing.T) {
	logger.Create.SetPath("./log/ek2.log").Register()
	logger.Info("123", `{"test":"1221"}`)
	logger.Default("./log/ek.log").Register()
	logger.Info("3333", `{"test":"22222"}`)
}
