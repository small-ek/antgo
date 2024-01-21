package test

import (
	"fmt"
	"github.com/small-ek/antgo/utils/conv"
	"github.com/small-ek/antgo/utils/plugins"
	"testing"
)

func TestPlugins(t *testing.T) {
	plugins.New().Register("example", &ExamplePlugin{Name: "123"})
	plugins.New().Uninstall("example")
	plugins.New().Register("example", &ExamplePlugin{Name: "1233"})
	plugins.New().Register("example2", &ExamplePlugin{Name: "456789"})
	for _, plugin := range plugins.List() {
		plugin.Before()
		plugin.After("1111")
	}
}

type ExamplePlugin struct {
	Name string
}

func (ep *ExamplePlugin) Before() interface{} {
	// 在操作之前执行的逻辑
	return "Before Logic"
}

func (ep *ExamplePlugin) After(data ...interface{}) {
	fmt.Println(ep.Name)
	fmt.Println(conv.String(data))
	// 在操作之后执行的逻辑
	fmt.Println("After Logic")
}
