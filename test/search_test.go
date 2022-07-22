package test

import (
	"fmt"
	"github.com/small-ek/antgo/container/array"
	"testing"
)

func TestSearch(t *testing.T) {
	s_f32 := []interface{}{1, 2, 3, 4}
	s_key := 1.0
	isEq := array.SearchInterface(s_f32, s_key)
	fmt.Println(isEq)
}
