package test

import (
	"github.com/small-ek/antgo/encoding/json"
	"log"
	"testing"
)

func TestJson(t *testing.T) {
	log.Println(string(json.Open("i18n.json")))
	for i := 0; i < 1; i++ {
		jsonStr := `[{"users" : {
	    "count" : 2,
	    "list"  : [
	        {"name" : "Ming", "score" : 60},
	        {"name" : "John", "score" : 99.5}
	    ]
		}
	}]`
		var result = json.Decode(jsonStr).Get(0).Get("users").Get("list").Array()
		log.Println(result)
		log.Println(json.Encode(map[string]string{"name": "21"}))

		var result2 = json.Decode(jsonStr).Read("0.users.list").Array()
		log.Println(result2)
	}
}
