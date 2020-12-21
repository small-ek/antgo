package test

import (
	"github.com/small-ek/antgo/os/logs"
	"testing"
)

func TestErr(t *testing.T) {
	test1()
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
