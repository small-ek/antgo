package test

import (
	"github.com/small-ek/ginp/conv"
	"github.com/small-ek/ginp/os/cache"
	"log"
	"testing"
)

func TestCache(t *testing.T) {
	var test string
	for i := 0; i < 100; i++ {
		/*var caches = cache.New{
			Key:   conv.String(i),
			Group: "test",
		}*/
		/*cache.Set("test"+conv.String(i), "testtes")*/
		var result=cache.GetOrSet("test22"+conv.String(i), "testtes")
		log.Println(string(result))
		var result2 = cache.GetOrSet("test22"+conv.String(i), "testtes")
		conv.BytesToStruct(result2, &test)

		log.Println(test)
	}

}
