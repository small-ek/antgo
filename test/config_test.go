package test

import (
	"fmt"
	"github.com/small-ek/antgo/os/config"
	"testing"
)

func TestConfig(t *testing.T) {
	err := config.New("./config.toml").Regiter()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(config.Get("zh"))
}
