package test

import (
	"github.com/small-ek/ginp/container"
	"github.com/small-ek/ginp/conv"
	"log"
	"testing"
)

func TestMap(t *testing.T) {
	var data = container.Map{Map: make(map[string]interface{})}

	for i := 0; i < 21; i++ {
		go func(n int) {
			data.Set("test"+conv.String(i), "ek"+conv.String(i))
			log.Println("test" + conv.String(i))
			log.Println(data.Get("test" + conv.String(i)))

		}(i)
	}
	var test="123"
	log.Println(test)
	log.Println(data)
}
