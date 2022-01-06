package test

import (
	"github.com/small-ek/antgo/utils/conv"
	"log"
	"testing"
)

func TestInts(t *testing.T) {
	var result = conv.Ints([]string{"1", "2", "3"})
	log.Println(result)
}
