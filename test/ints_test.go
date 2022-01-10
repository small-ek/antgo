package test

import (
	"github.com/small-ek/antgo/utils/conv"
	"log"
	"testing"
)

func TestInts(t *testing.T) {
	var str = conv.Ints([]string{"1", "2", "3"})
	log.Println(str)
	var str2 = conv.Ints("[1,2,3]")
	log.Println(str2)
}
