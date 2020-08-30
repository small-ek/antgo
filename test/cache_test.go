package test

import (
	"github.com/small-ek/ginp/conv"
	"github.com/small-ek/ginp/os/cache"
	"log"
	"testing"
)

func TestCache(t *testing.T) {
	var ek string

	for i := 0; i < 100; i++ {
		cache.Set("ek"+conv.String(i), "ek"+conv.String(i))
		var result = cache.Get("ek" + conv.String(i))
		conv.BytesToStruct(result, &ek)
		log.Println(ek)
	}
}
