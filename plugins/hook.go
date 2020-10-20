package plugins

//pluginFunc Define an interface, there are two methods
type pluginFunc interface {
	Before() interface{} //Before what
	After()              //After what
}

//New Define a class to store our plugin
type New struct {
	List map[string]pluginFunc //Plugin list
}

//Init Initialize the plugin
func (p *New) Init() {
	p.List = make(map[string]pluginFunc)
}

// Register plugin
func (p *New) Register(path, name string, plugin pluginFunc) {
	if path == name {
		p.List[name] = plugin
	}
}
