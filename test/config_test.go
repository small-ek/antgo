package test

import (
	"fmt"
	"github.com/small-ek/antgo/os/config"
	"testing"
)

func TestConfig(t *testing.T) {
	err := config.New().SetType("toml").AddRemoteProvider("etcd3", "http://127.0.0.1:2379", "/local/common.toml")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(config.Get("app"))
}
