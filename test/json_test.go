package test

import (
	"github.com/small-ek/ginp/encoding/json"
	"log"
	"testing"
)

func TestJson(t *testing.T) {
	//var data = `{
	//"users" : {
	//    "count" : 2,
	//    "list"  : [
	//        {"name" : "Ming", "score" : 60},
	//        {"name" : "John", "score" : 99.5}
	//    ]
	//	}
	//}`
	/*var bindRow = make(map[string]interface{})
	conv.Struct(&bindRow, data)
	log.Println(bindRow)*/
	jsonStr := `
        [{
	"users" : {
	    "count" : 2,
	    "list"  : [
	        {"name" : "Ming", "score" : 60},
	        {"name" : "John", "score" : 99.5}
	    ]
		}
	}]
        `
	var result = json.DecodeArray(jsonStr).Array()
	log.Println(result)
}
