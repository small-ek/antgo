package test

import (
	"github.com/small-ek/antgo/crypto/hash"
	"log"
	"testing"
)

func TestHash(t *testing.T) {
	log.Println(hash.Md5("test"))
	log.Println(hash.Sha1("test"))
	log.Println(hash.Sha256("test"))
	log.Println(hash.Crc32("test"))
}
