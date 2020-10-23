package test

import (
	"github.com/small-ek/ginp/encoding/json"
	"log"
	"testing"
)

func TestJson(t *testing.T) {

	for i := 0; i < 1000; i++ {
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
	}
}
