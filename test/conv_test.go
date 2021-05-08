package test

import (
	"github.com/small-ek/antgo/conv"
	"log"
	"testing"
)

func TestConv(t *testing.T) {
	str1 := "1234A"
	log.Println(conv.Int(str1))
	str2 := "1234"
	log.Println(conv.Int(str2))
	log.Println(conv.Uint(str1))
	log.Println(conv.Uint(str1))
}
