package test

import (
	"github.com/small-ek/antgo/conv"
	"log"
	"testing"
)

func TestSlice(t *testing.T) {
	var result = conv.Strings([]int{1, 2, 3, 4})
	log.Println(result)

}
