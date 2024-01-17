package test

import (
	"github.com/small-ek/antgo/utils/ants"
	"log"
	"testing"
)

func TestPool(t *testing.T) {
	ants.InitPool(10000)
	err := ants.NewPool.Submit(func() {
		log.Println("12122112")
	})
	if err != nil {
		return
	}
}
