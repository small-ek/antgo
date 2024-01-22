package test

import (
	"github.com/small-ek/antgo/utils/conv"
	"log"
	"testing"
)

func TestConv(t *testing.T) {
	str1 := "1234A"
	log.Println(conv.Int(str1))
	str2 := "1234"
	log.Println(conv.Int(str2))
	log.Println(conv.Uint64(str2))
	var srt3 int = 3
	log.Println(conv.Uint64(srt3))
	var test float32 = 12
	log.Println(conv.Bytes(test))
}

type JsonStr struct {
	Name string
}

func TestConvMap(t *testing.T) {
	mapStr := `{"name":"Hello"}`
	log.Println(conv.Map(mapStr))
	MapJson := []JsonStr{{Name: "Hello2"}}
	log.Println(conv.Maps(MapJson))
}

func BenchmarkConv(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) { //并发
		for pb.Next() {
			str1 := "1234A"
			log.Println(conv.Int(str1))
			str2 := "1234"
			log.Println(conv.Int(str2))
			log.Println(conv.Uint64(str2))
			var srt3 int = 3
			log.Println(conv.Uint64(srt3))
			var test float32 = 12
			log.Println(conv.Bytes(test))
		}
	})
}
