package test

import (
	"github.com/small-ek/antgo/container/maps"
	"github.com/small-ek/antgo/conv"
	"log"
	"sync"
	"testing"
)

func TestMap(t *testing.T) {

	var data = maps.New()
	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func(n int) {
			var key = conv.String(n)
			var keys = "test" + key
			var value = "ek" + key
			log.Println(keys)
			log.Println(value)
			data.Set(keys, value)

			/*t.Logf("k=:%v,v:=%v\n,c=:%c", key, data.Get(keys), data.Count())*/
			wg.Done()
		}(i)
	}
	wg.Wait()
	var test = "123"
	log.Println(test)
	/*	log.Println(data)*/
}
