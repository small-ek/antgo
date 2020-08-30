package test

import (
	"github.com/small-ek/ginp/conv"
	"github.com/small-ek/ginp/os/storage"
	"testing"
)

func TestLeveldb(t *testing.T) {
	/*var ek string*/
	var s = storage.Open("./tmp")

	for i := 0; i < 100; i++ {
		s.Set("ek"+conv.String(i), "test"+conv.String(i))
		  
	}
	defer s.Db.Close()
}
