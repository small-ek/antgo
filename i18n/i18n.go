package i18n

import (
	"github.com/small-ek/ginp/os/config"
)

//New i18n
type New struct {
	Tag string
	Msg map[string]interface{}
}

var Result *New

//SetPath Set path
func SetPath(path string) {
	config.SetPath(path)
}

//SetLanguage Set language
func SetLanguage(languages string) {
	var getLang = config.Default().Get(languages).Map()
	var child = make(map[string]interface{})

	for key, value := range getLang {
		if languages != "" && key != "" && value != nil {
			child[key] = value
		}
	}

	Result = &New{
		Tag: languages,
		Msg: child,
	}
}

//Get language
func Get(name string) string {
	return Result.Msg[name].(string)
}
