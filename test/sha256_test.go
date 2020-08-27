package test

import (
	"github.com/small-ek/ginp/crypto/sha256"
	"log"
	"testing"
)

func TestSha256(t *testing.T)  {
	for i:=0;i<100;i++{
		var str1= sha256.Create("admin");
		log.Println(str1)
	}
}