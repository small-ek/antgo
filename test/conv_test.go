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

	log.Println(conv.Float64(str1))
	log.Println(conv.Float32(str2))
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

func TestInterfaces(t *testing.T) {
	mapStr := `[1,2,3,4]`
	log.Println(conv.Interfaces(mapStr))
	MapJson := []string{"1", "2", "3", "4"}
	log.Println(conv.Interfaces(MapJson))
}

type Test struct {
	Name  string `json:"name"`
	Name2 string `json:"name2"`
}

type Test2 struct {
	Name string `json:"name"`
}

func TestStruct(t *testing.T) {
	test := Test{}
	test2 := Test2{Name: "222222222222"}
	conv.ToStruct(test2, &test)
	log.Println(test)
}

type Test3 struct {
	Name  string `json:"name"`
	Name2 string `json:"name2"`
}

func TestJSON(t *testing.T) {
	test3 := Test3{}
	str := `{"name":"123"}`
	conv.UnmarshalJSON([]byte(str), &test3)
	log.Println(conv.String(test3))

}

type User struct {
	ID   int
	Name string
	Age  int
}

func TestGob(t *testing.T) {
	user := User{
		ID:   1,
		Name: "Alice",
		Age:  30,
	}

	// 将 User 数据结构转换为字节切片
	dataBytes, err := conv.GobEncoder(user)
	if err != nil {
		log.Fatalf("Error while converting to bytes: %v", err)
	}
	log.Println(dataBytes)
	user2 := User{}
	conv.GobDecoder(dataBytes, &user2)
	log.Println(user2)
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
