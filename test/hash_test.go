package test

import (
	"github.com/small-ek/antgo/crypto/ahash"
	"log"
	"testing"
)

func TestHash(t *testing.T) {
	log.Println(ahash.MD5("test"))
	log.Println(ahash.SHA1("test"))
	log.Println(ahash.SHA256("test"))
	log.Println(ahash.SHA512("test"))
	log.Println(ahash.Crc32("test"))
}
