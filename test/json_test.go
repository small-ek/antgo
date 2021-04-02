package test

import (
	"github.com/small-ek/antgo/encoding/jsons"
	"log"
	"testing"
)

func TestJson(t *testing.T) {
	log.Println(string(jsons.Open("i18n.json")))
	for i := 0; i < 1; i++ {
		jsonStr := `[{"users" : {
	    "count" : 2,
	    "list"  : [
	        {"name" : "Ming", "score" : 60},
	        {"name" : "John", "score" : 99.5}
	    ]
		}
	}]`
		var result = jsons.Decode(jsonStr).Get(0).Get("users").Get("list").Array()
		log.Println(result)
		log.Println(jsons.Encode(map[string]string{"name": "21"}))

		var result2 = jsons.Decode(jsonStr).Read("0.users.list").Array()
		log.Println(result2)
	}
}
