package test

import (
	"flag"
	"github.com/small-ek/antgo/i18n"
	"log"
	"reflect"
	"testing"
)

func TestI18n(t *testing.T) {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	flag.Parse()

	i18n.SetPath("i18n.json")
	for i := 0; i < 1000; i++ {
		i18n.SetLanguage("zh")
		var result = i18n.Get("hello")
		log.Println(result)
		log.Println(reflect.TypeOf(result))
	}
}
