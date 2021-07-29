package config

import (
	"github.com/BurntSushi/toml"
	"github.com/small-ek/antgo/conv"
	"github.com/small-ek/antgo/os/logs"
	"os"
	"strings"
)

//Data config data
var Data map[string]interface{}

//SetPath Set path.
func SetPath(path string) {
	if _, err := toml.DecodeFile(path, &Data); err != nil {
		logs.Error(err.Error())
		os.Exit(1)
	}
}

//Result ...
type Config struct {
	Child interface{}
}

//Decode config
func Decode() *Config {
	return &Config{
		Child: Data,
	}
}

//Get config
func (c *Config) Get(name interface{}) *Config {
	var child = c.Child
	switch child.(type) {
	case map[string]interface{}:
		return &Config{
			Child: child.(map[string]interface{})[conv.String(name)],
		}
	case map[string]string:
		return &Config{
			Child: child.(map[string]string)[conv.String(name)],
		}
	case map[string]int:
		return &Config{
			Child: child.(map[string]string)[conv.String(name)],
		}
	case []interface{}:
		return &Config{
			Child: child.([]interface{})[conv.Int(name)],
		}
	case []string:
		return &Config{
			Child: child.([]interface{})[conv.Int(name)],
		}
	case []int:
		return &Config{
			Child: child.([]interface{})[conv.Int(name)],
		}
	case []int64:
		return &Config{
			Child: child.([]interface{})[conv.Int(name)],
		}
	}
	return &Config{
		Child: child,
	}
}

//Read Parse config according to point split
func (c *Config) Read(name string) *Config {
	var list = strings.Split(name, ".")
	for i := 0; i < len(list); i++ {
		var value = list[i]
		var result = c.Get(value)
		c.Child = result.Child
	}
	return c
}

//String Data type conversion.
func (c *Config) String() string {
	if c.Child == nil {
		return ""
	}
	return conv.String(c.Child)
}

//Strings Data type conversion.
func (c *Config) Strings() []string {
	return conv.Strings(c.Child)
}

//Byte Data type conversion.
func (c *Config) Byte() byte {
	return conv.Byte(c.Child)
}

//Bytes Data type conversion.
func (c *Config) Bytes() []byte {
	return conv.Bytes(c.Child)
}

//Int Data type conversion.
func (c *Config) Int() int {
	if c.Child == nil {
		return 0
	}
	return conv.Int(c.Child)
}

//Bool Data type conversion.
func (c *Config) Bool() bool {
	return conv.Bool(c.Child)
}

//Ints Data type conversion.
func (c *Config) Ints() []int {
	return conv.Ints(c.Child)
}

//Int64 Data type conversion.
func (c *Config) Int64() int64 {
	if c.Child == nil {
		return 0
	}
	return conv.Int64(c.Child)
}

//Float64 Data type conversion.
func (c *Config) Float64() float64 {
	if c.Child == nil {
		return 0
	}
	return conv.Float64(c.Child)
}

//Map Data type conversion.
func (c *Config) Map() map[string]interface{} {
	if c.Child == nil {
		return nil
	}
	return c.Child.(map[string]interface{})
}

//Array Data type conversion.
func (c *Config) Array() []interface{} {
	if c.Child == nil {
		return nil
	}
	return c.Child.([]interface{})
}

//Uint Data type conversion.
func (c *Config) Uint() uint {
	return conv.Uint(c.Child)
}

//Uint Data type conversion.
func (c *Config) Uint8() uint8 {
	return conv.Uint8(c.Child)
}

//Uint Data type conversion.
func (c *Config) Uint16() uint16 {
	return conv.Uint16(c.Child)
}

//Uint Data type conversion.
func (c *Config) Uint32() uint32 {
	return conv.Uint32(c.Child)
}

//Uint Data type conversion.
func (c *Config) Uint64() uint64 {
	return conv.Uint64(c.Child)
}
