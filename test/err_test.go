package test

import (
	"github.com/small-ek/antgo/os/logs"
	"log"
	"testing"
)

func TestErr(t *testing.T) {
	for i := 0; i < 20; i++ {
		test1()
		log.Println("21111")
	}
}

func test1() {
	test2()
}

func test2() {
	test3()
}

func test3() {
	logs.Error("111")
}
