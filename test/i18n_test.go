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

	lang :=i18n.New("./language/i18n", "zh-CN.toml")
	var result = lang.Get("hello")
	log.Println(result)
	log.Println(reflect.TypeOf(result))
}
