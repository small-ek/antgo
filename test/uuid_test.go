package test

import (
	"log"
	"testing"
)

func TestUuid(t *testing.T) {
	for i := 0; i < 100; i++ {
		uuid1 := auuid.Create()
		log.Println(uuid1)
		uuid2 := auuid.NewDCEGroup()
		log.Println(uuid2)
		uuid3 := auuid.NewDCEPerson()
		log.Println(uuid3)
		uuid4 := auuid.NewRandom()
		log.Println(uuid4)
	}
}
