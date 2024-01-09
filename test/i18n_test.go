package test

import (
	"flag"
	"github.com/small-ek/antgo/i18n"
	"log"
	"testing"
)

func TestI18n(t *testing.T) {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	flag.Parse()

	lang := i18n.New("./language/i18n", "zh-CN")
	str := lang.T("hello", "12", "23")
	log.Println(str)

	str2 := lang.TOption("hello", "en")
	log.Println(str2)

	str3 := lang.T("hello", "12", "23")
	log.Println(str3)

	str4 := lang.T("common.name", "2222")
	log.Println(str4)
}
