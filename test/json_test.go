package test

import (
	"github.com/small-ek/ginp/encoding/jsons"
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
        {
	"users" : {
	    "count" : 2,
	    "list"  : [
	        {"name" : "Ming", "score" : 60},
	        {"name" : "John", "score" : 99.5}
	    ]
		}
	}
        `
	log.Println(jsonStr)
	/*var result = jsons.DecodeJson(jsonStr).Get("users")
	log.Println(result["12"])*/
}
