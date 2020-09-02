package test

import (
	"github.com/small-ek/ginp/os/logs"
	"testing"
)

func TestErr(t *testing.T) {
	for i := 0; i < 100; i++ {
		logs.Error("1234")
	}
}
