package test

import (
	"github.com/small-ek/ginp/crypto/md5"
	"log"
	"testing"
)

func TestMd5(t *testing.T)  {
	for i:=0;i<100;i++{
		log.Println(md5.Create("test"))
	}
}