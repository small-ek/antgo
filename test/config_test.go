package test

import (
	"github.com/small-ek/ginp/os/config"
	"log"
	"testing"
)

func TestConfig(t *testing.T) {
	config.SetPath("config.toml")
	/*log.Println(config.Get["Master"])*/
	var result = config.Decode().Get("en").Map()
	log.Println(result)
}
