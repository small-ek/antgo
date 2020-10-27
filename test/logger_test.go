package test

import (
	"github.com/small-ek/ginp/os/logger"
	"log"
	"runtime"
	"testing"
)

func TestLogger(t *testing.T) {
	logger.Default("./log/ek2.log").Register()
	logger.Info("123", `{"test":"1221"}`)
	logger.Write.Info("22333")
	test()
}
func test() {
	test23()
}
func test23() {
	pc, file, line, ok := runtime.Caller(2)
	log.Println(pc)
	log.Println(file)
	log.Println(line)
	log.Println(ok)
	f := runtime.FuncForPC(pc)
	log.Println(f.Name())

	pc, file, line, ok = runtime.Caller(0)
	log.Println(pc)
	log.Println(file)
	log.Println(line)
	log.Println(ok)
	f = runtime.FuncForPC(pc)
	log.Println(f.Name())

	pc, file, line, ok = runtime.Caller(1)
	log.Println(pc)
	log.Println(file)
	log.Println(line)
	log.Println(ok)
	f = runtime.FuncForPC(pc)
	log.Println(f.Name())
}
