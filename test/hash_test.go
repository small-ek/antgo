package test

import (
	"github.com/small-ek/antgo/crypto/ahash"
	"log"
	"testing"
)

func TestHash(t *testing.T) {
	log.Println(ahash.Md5("test"))
	log.Println(ahash.Sha1("test"))
	log.Println(ahash.Sha256("test"))
	log.Println(ahash.Crc32("test"))
}
