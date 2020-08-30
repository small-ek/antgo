package test

import (
	"github.com/small-ek/ginp/conv"
	"github.com/small-ek/ginp/os/storage"
	"log"
	"testing"
)

func TestLeveldb(t *testing.T) {
	/*var ek string*/
	var s = storage.Open("./tmp")
	defer s.Db.Close()
	for i := 0; i < 200; i++ {
		/*s.Set("ek"+conv.String(i), "test"+conv.String(i))*/
		var test, err = s.Get("ek" + conv.String(i))
		log.Println(test)
		log.Println(err)
	}

}
