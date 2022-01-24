package test

import (
	"github.com/small-ek/antgo/utils/str"
	"log"
	"testing"
)

func TestStr(t *testing.T) {
	str1 := `"123"`
	log.Println(str.ClearQuotes(str1))
}
