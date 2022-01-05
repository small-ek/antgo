package i18n

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/small-ek/antgo/conv"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
)

type I18n struct {
	Path     string
	Language string
	Data     map[string]interface{}
	Source   map[string]interface{}
}

//Datas ...

//Maps ...
var Maps map[string]interface{}

//Array ...
var Array []interface{}

//SetPath Set path
func New(filePath, Language string) *I18n {
	data := make(map[string]interface{})
	filePath = fmt.Sprintf("%s/%s", filePath, Language)
	fileNameWithSuffix := path.Base(filePath)
	fileType := path.Ext(fileNameWithSuffix)
	switch fileType {
	case ".toml":
		if _, err := toml.DecodeFile(filePath, &data); err != nil {
			panic(err.Error())
		}
		break
	case ".yml", ".yaml":
		file, _ := os.Open(filePath) //test.yaml由下一个例子生成
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
		Path:     filePath,
		Language: Language,
		Source:   data,
		Data:     make(map[string]interface{}),
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

//Get language
func (i *I18n) T(key string, args ...interface{}) string {
	format := key

	if _, ok := i.Data[key]; ok {
		format = conv.String(i.Data[key])
	} else {
		for _, row := range i.Source {
			if row[0] == key {
				i.Data[key] = row[1]
				format = row[1]
				break
			}
		}
	}
	format = i.preArgs(format, args...)
	return format
}

func (i *I18n) preArgs(format string, args ...interface{}) string {
	if len(args) > 0 {
		format = fmt.Sprintf(format, args...)
	}
	return format
}
