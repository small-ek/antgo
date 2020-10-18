package i18n

import (
	"github.com/small-ek/ginp/os/config"
)

type Data struct {
	Tag string
	Msg map[string]interface{}
}

var Result *Data

func SetPath(path string) {
	config.SetPath(path)
}

func SetLanguage(languages string) {
	var getLang = config.Default().Get(languages).Map()
	var child = make(map[string]interface{})

	for key, value := range getLang {
		if languages != "" && key != "" && value != nil {
			child[key] = value
		}
	}
	Result = &Data{
		Tag: languages,
		Msg: child,
	}
}

func Get(name string) string {
	return Result.Msg[name].(string)
}
