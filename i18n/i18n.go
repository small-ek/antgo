package i18n

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/small-ek/antgo/conv"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type I18n struct {
	Path     string
	Language string
	Data     map[string]string
	Source   map[string]string
}

//SetPath Set path
func New(prefixPath, language string, defaultType ...string) *I18n {
	data := make(map[string]string)

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
		Data:     make(map[string]string),
	}
}

//SetPath language setting
func (i *I18n) SetPath(path string) {
	i.Path = path
}

//SetLanguage language setting
func (i *I18n) SetLanguage(language string) {
	i.Language = language
}

//T Get language
func (i *I18n) T(key string, args ...interface{}) string {
	format := key

	if _, ok := i.Data[key]; ok {
		format = conv.String(i.Data[key])
	} else {
		for value, row := range i.Source {
			i.Data[key] = row
			if value == key {
				format = conv.String(row)
				break
			}
		}
	}
	format = i.preArgs(format, args...)
	return format
}

//TOption Choose language translation
func (i *I18n) TOption(key string, language string, args ...interface{}) string {
	lang := New(i.Path, language)

	format := key
	for value, row := range lang.Source {
		lang.Data[key] = row
		if value == key {
			format = conv.String(row)
			break
		}
	}

	format = lang.preArgs(format, args...)
	return format
}

//preArgs ...
func (i *I18n) preArgs(format string, args ...interface{}) string {
	if len(args) > 0 {
		format = fmt.Sprintf(format, args...)
	}
	return format
}
