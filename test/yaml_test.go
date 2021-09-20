package test

import (
	"github.com/small-ek/antgo/encoding/ayaml"
	"log"
	"testing"
)

var yamlStr string = `
#即表示url属性值；
url: https://github.com/small-ek/antgo

#数组，即表示server为[a,b,c]
server:
    - 192.168.1.1
    - 192.168.1.2
#常量
pi: 3.14   #定义一个数值3.14
hasChild: true  #定义一个boolean值
name: '你好YAML'   #定义一个字符串
`
var yamlErr string = `
{"name":"123"}
`

func TestYaml(t *testing.T) {
	result, err := ayaml.Decode([]byte(yamlStr))
	if err != nil {
		t.Errorf("encode failed. %v", err)
		return
	}
	log.Println(result)

	result2 := make(map[string]interface{})
	err = ayaml.DecodeTo([]byte(yamlStr), &result2)
	log.Println(result2)

	m := make(map[string]string)
	m["yaml"] = yamlStr
	res, err := ayaml.Encode(m)
	if err != nil {
		t.Errorf("encode failed. %v", err)
		return
	}
	log.Println(string(res))

	jsonyaml, err := ayaml.ToJson([]byte(yamlErr))
	if err != nil {
		t.Errorf("ToJson failed. %v", err)
		return
	}
	log.Println(string(jsonyaml))
}
