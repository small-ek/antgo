package test

import (
	"github.com/small-ek/antgo/utils/conv"
	"log"
	"testing"
)

func TestSlice(t *testing.T) {
	var result = conv.Strings([]int{1, 2, 3, 4})
	log.Println(result)
	var result2 = conv.Strings(`["1", "2", "3", "4"]`)
	log.Println(result2)
	var result3 = conv.Strings([]string{"1", "2"})
	log.Println(result3)
	var result4 = conv.Strings([]int64{1, 2, 3, 4})
	log.Println(result4)
}
