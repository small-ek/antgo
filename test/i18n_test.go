package test

import (
	"github.com/small-ek/ginp/i18n"
	"log"
	"testing"
)

/*var Lang map[string]map[string]interface{}
var Tag string*/

func TestI18n(t *testing.T) {
	i18n.SetPath("config.toml")
	i18n.SetLanguage("en")
	log.Println(i18n.Get("hello"))
}
