package config

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/small-ek/antgo/conv"
	"github.com/small-ek/antgo/os/logs"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

//Data config data
var Data map[string]interface{}

//SetPath Set path.
func SetPath(filePath string) {
	fileNameWithSuffix := path.Base(filePath)
	fileType := path.Ext(fileNameWithSuffix)
	switch fileType {
	case ".toml":
		if _, err := toml.DecodeFile(filePath, &Data); err != nil {
			panic(err.Error())
		}
		break
	case ".yml", ".yaml":
		file, _ := os.Open(filePath) //test.yaml由下一个例子生成
		defer file.Close()
		ydecode := yaml.NewDecoder(file)
		if err := ydecode.Decode(&Data); err != nil {
			panic(err.Error())
		}
		break
	case "json":
		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			logs.Error(err.Error())
			return
		}
		err2 := json.Unmarshal(bytes, &Data)
		if err2 != nil {
			panic(err.Error())
			return
		}
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

//Next config
func (c *Config) Next(name interface{}) *Config {
	var child = c.Child
	switch child.(type) {
	case map[string]interface{}:
		return &Config{
			Child: child.(map[string]interface{})[conv.String(name)],
		}
	case map[interface{}]interface{}:
		return &Config{
			Child: child.(map[interface{}]interface{})[conv.String(name)],
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

//Get Parse config according to point split
func (c *Config) Get(name string) *Config {
	var list = strings.Split(name, ".")
	for i := 0; i < len(list); i++ {
		var value = list[i]
		var result = c.Next(value)
		c.Child = result.Child
	}
	return c
}

//String Data type conversion.
func (c *Config) String() string {
	if c.Child == nil {
		return ""
	}
	defer c.End()
	return conv.String(c.Child)
}

//End Data type conversion.
func (c *Config) End() {
	c.Child = Data
}

//Strings Data type conversion.
func (c *Config) Strings() []string {
	defer c.End()
	return conv.Strings(c.Child)
}

//Byte Data type conversion.
func (c *Config) Byte() byte {
	defer c.End()
	return conv.Byte(c.Child)
}

//Bytes Data type conversion.
func (c *Config) Bytes() []byte {
	defer c.End()
	return conv.Bytes(c.Child)
}

//Int Data type conversion.
func (c *Config) Int() int {
	if c.Child == nil {
		return 0
	}
	defer c.End()
	return conv.Int(c.Child)
}

//Bool Data type conversion.
func (c *Config) Bool() bool {
	defer c.End()
	return conv.Bool(c.Child)
}

//Ints Data type conversion.
func (c *Config) Ints() []int {
	defer c.End()
	return conv.Ints(c.Child)
}

//Int64 Data type conversion.
func (c *Config) Int64() int64 {
	if c.Child == nil {
		return 0
	}
	defer c.End()
	return conv.Int64(c.Child)
}

//Float64 Data type conversion.
func (c *Config) Float64() float64 {
	if c.Child == nil {
		return 0
	}
	defer c.End()
	return conv.Float64(c.Child)
}

//Map Data type conversion.
func (c *Config) Map() map[string]interface{} {
	if c.Child == nil {
		return nil
	}
	defer c.End()
	return c.Child.(map[string]interface{})
}

//Maps Data type conversion.
func (c *Config) Maps() []map[string]interface{} {
	if c.Child == nil {
		return nil
	}
	defer c.End()
	return c.Child.([]map[string]interface{})
}

//Array Data type conversion.
func (c *Config) Array() []interface{} {
	if c.Child == nil {
		return nil
	}
	defer c.End()
	return c.Child.([]interface{})
}

//Uint Data type conversion.
func (c *Config) Uint() uint {
	defer c.End()
	return conv.Uint(c.Child)
}

//Uint Data type conversion.
func (c *Config) Uint8() uint8 {
	defer c.End()
	return conv.Uint8(c.Child)
}

//Uint Data type conversion.
func (c *Config) Uint16() uint16 {
	defer c.End()
	return conv.Uint16(c.Child)
}

//Uint Data type conversion.
func (c *Config) Uint32() uint32 {
	defer c.End()
	return conv.Uint32(c.Child)
}

//Uint Data type conversion.
func (c *Config) Uint64() uint64 {
	defer c.End()
	return conv.Uint64(c.Child)
}

//Interfaces Data type conversion.
func (c *Config) Interfaces() []interface{} {
	return conv.Interfaces(c.Child)
}

//Interface Data type conversion.
func (c *Config) Interface() interface{} {
	return c.Child
}
