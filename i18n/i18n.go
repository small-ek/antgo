package i18n

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/small-ek/antgo/utils/conv"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

//I18n
type I18n struct {
	Path     string
	Language string
	Type     string
	Data     map[string]string
	Source   map[string]string
}

//New initialization
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

//T language translation
func (i18n *I18n) T(key string, args ...interface{}) string {
	format := key

	if _, ok := i18n.Data[key]; ok {
		format = conv.String(i18n.Data[key])
	}
	format = i18n.preArgs(format, args...)
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
