package test

import (
	"fmt"
	"github.com/small-ek/antgo/encoding/abinary"
	"testing"
)

func TestBinary(t *testing.T) {
	//序列化
	var dataA uint64 = 6010
	buffer := abinary.Encode(dataA)

	byteA := buffer.Bytes()
	fmt.Println("序列化后：", byteA)

	//反序列化
	var res uint64
	abinary.Decode(byteA, &res)

	fmt.Println("反序列化后：", res)
}
