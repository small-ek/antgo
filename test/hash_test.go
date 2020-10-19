package test

import (
	"github.com/small-ek/ginp/crypto/hash"
	"log"
	"testing"
)

func TestHash(t *testing.T) {
	log.Println(hash.MD5("test"))
	log.Println(hash.Sha("test"))
	log.Println(hash.Sha256("test"))
}
