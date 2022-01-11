package test

import (
	"github.com/small-ek/antgo/os/config"
	"log"
	"testing"
)

func TestConfig(t *testing.T) {
	config.SetPath("config.yaml")
	/*log.Println(config.Get["Master"])*/
	var result = config.Decode().Get("en").Get("hello").String()
	log.Println(result)
	config.SetPath("config2.toml")
	var result2 = config.Decode().Get("en").Get("hello2").String()
	log.Println(result2)
	var result3 = config.Decode().Get("en.hello2").String()
	log.Println(result3)
}
