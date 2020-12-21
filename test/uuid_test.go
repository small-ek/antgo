package test

import (
	"github.com/small-ek/antgo/crypto/uuid"
	"log"
	"testing"
)

func TestUuid(t *testing.T) {
	for i := 0; i < 100; i++ {
		uuid1 := uuid.Create()
		log.Println(uuid1)
		uuid2 := uuid.NewDCEGroup()
		log.Println(uuid2)
		uuid3 := uuid.NewDCEPerson()
		log.Println(uuid3)
		uuid4 := uuid.NewRandom()
		log.Println(uuid4)
	}
}
