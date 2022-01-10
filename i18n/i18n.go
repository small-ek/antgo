package i18n

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/small-ek/antgo/utils/conv"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

//I18n
type I18n struct {
	Path     string
	Language string
	Type     string
	Data     map[string]interface{}
	Source   map[string]interface{}
	Child    interface{}
}

//New initialization
func New(prefixPath, language string, defaultType ...string) *I18n {
	data := make(map[string]interface{})

	types := "toml"
	if len(defaultType) > 0 {
		types = defaultType[0]
	}

	filePath := fmt.Sprintf("%s/%s.%s", prefixPath, language, types)

	switch types {
	case "toml":
		if _, err := toml.DecodeFile(filePath, &data); err != nil {
			panic(err.Error())
		}
		break
	case "yml", "yaml":
		file, _ := os.Open(filePath)
		defer file.Close()
		ydecode := yaml.NewDecoder(file)
		if err := ydecode.Decode(&data); err != nil {
			panic(err.Error())
		}
		break
	case "json":
		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(err.Error())
		}
		err2 := json.Unmarshal(bytes, &data)
		if err2 != nil {
			panic(err.Error())
		}
	}
	return &I18n{
		Path:     prefixPath,
		Language: language,
		Source:   data,
		Type:     types,
		Data:     data,
	}
}

//SetPath Language pack path
func (i *I18n) SetPath(path string) {
	i.Path = path
}

//SetLanguage language setting
func (i *I18n) SetLanguage(language string) {
	i.Language = language
}

//Next config
func (i18n *I18n) Next(name interface{}) *I18n {
	var child = i18n.Child
	switch child.(type) {
	case map[string]interface{}:
		return &I18n{
			Child: child.(map[string]interface{})[conv.String(name)],
		}
	case map[interface{}]interface{}:
		return &I18n{
			Child: child.(map[interface{}]interface{})[conv.String(name)],
		}
	case map[string]string:
		return &I18n{
			Child: child.(map[string]string)[conv.String(name)],
		}
	case map[string]int:
		return &I18n{
			Child: child.(map[string]string)[conv.String(name)],
		}
	case []interface{}:
		return &I18n{
			Child: child.([]interface{})[conv.Int(name)],
		}
	case []string:
		return &I18n{
			Child: child.([]interface{})[conv.Int(name)],
		}
	case []int:
		return &I18n{
			Child: child.([]interface{})[conv.Int(name)],
		}
	case []int64:
		return &I18n{
			Child: child.([]interface{})[conv.Int(name)],
		}
	case string:
		return &I18n{
			Child: i18n.Child,
		}
	}
	return &I18n{
		Child: child,
	}
}

//T language translation
func (i18n *I18n) T(key string, args ...interface{}) string {
	list := strings.Split(key, ".")
	format := key

	if len(list) > 1 {
		i18n.Child = i18n.Data
		for i := 0; i < len(list); i++ {
			row := list[i]
			next := i18n.Next(row)
			i18n.Child = next.Child
		}
		format = conv.String(i18n.Child)
	}

	if len(list) <= 1 {
		if _, ok := i18n.Data[key]; ok {
			format = conv.String(i18n.Data[key])
		}
	}
	format = i18n.preArgs(format, args...)

	if format == "" {
		return key
	}
	return format
}

//TOption Choose language translation
func (i18n *I18n) TOption(key string, language string, args ...interface{}) string {
	lang := New(i18n.Path, language)

	format := key
	if _, ok := lang.Data[key]; ok {
		format = conv.String(lang.Data[key])
	}

	format = lang.preArgs(format, args...)
	return format
}

//preArgs Formatted text
func (i *I18n) preArgs(format string, args ...interface{}) string {
	if len(args) > 0 {
		format = fmt.Sprintf(format, args...)
	}
	return format
}
