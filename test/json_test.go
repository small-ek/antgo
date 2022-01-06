package test

import (
	"github.com/small-ek/antgo/encoding/ajson"
	"log"
	"testing"
)

func TestJson(t *testing.T) {
	//log.Println(string(ajson.Open("i18n.json")))
	for i := 0; i < 1; i++ {
		jsonStr := `[{"users" : {
	    "count" : 2,
	    "list"  : [
	        {"name" : "Ming", "score" : 60},
	        {"name" : "John", "score" : 99.5}
	    ]
		}
	}]`
		var result = ajson.Decode([]byte(jsonStr)).Get("0").Get("users").Get("list").Array()
		log.Println(result)
		log.Println(ajson.Encode(map[string]interface{}{"name": "21"}))

		var result2 = ajson.Decode([]byte(jsonStr)).Get("0.users.list").Interface()
		log.Println(result2)
	}
}
