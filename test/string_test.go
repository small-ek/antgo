package test

import (
	"github.com/small-ek/antgo/utils/conv"
	"log"
	"testing"
)

func TestString(t *testing.T) {
	for i := 0; i < 100; i++ {
		log.Println(conv.String(123))
		log.Println(conv.String(123.123))
		log.Println(conv.String("123.123"))
	}
}
