package plugins

import "sync"

// pluginFunc Define an interface, there are two methods
type pluginFunc interface {
	Before() interface{}       //Before what
	After(data ...interface{}) //After what
}

// New Define a class to store our plugin
type newPlugins struct {
	List map[string]pluginFunc //Plugin list
}

var once sync.Once
var Plugins *newPlugins

func New() *newPlugins {
	once.Do(func() {
		Plugins = &newPlugins{List: map[string]pluginFunc{}}
	})
	return Plugins
}

// List
func List() map[string]pluginFunc {
	return Plugins.List
}

// Register plugin
func (p *newPlugins) Register(name string, plugin pluginFunc) {
	_, ok := p.List[name]
	if !ok {
		p.List[name] = plugin
	}
}

// Uninstall plugin
func (p *newPlugins) Uninstall(name string) {
	_, ok := p.List[name]
	if ok {
		delete(p.List, name)
	}
}
