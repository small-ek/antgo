package aini

import (
	"log"
	"testing"
)

var iniStr string = `

;注释
aa=bb
[addr] 
#注释
ip = 127.0.0.1
port=9001
enable=true

	[DBINFO]
	type=mysql
	user=root
	password=password
[键]
呵呵=值

`

func TestIni(t *testing.T) {
	result, err := Decode([]byte(iniStr))
	if err != nil {
		t.Errorf("encode failed. %v", err)
		return
	}
	log.Println(result)

	res, err := Encode(result)
	if err != nil {
		t.Errorf("encode failed. %v", err)
		return
	}
	log.Println(string(res))

	jsonyaml, err := ToJson([]byte(iniStr))
	if err != nil {
		t.Errorf("ToJson failed. %v", err)
		return
	}
	log.Println(string(jsonyaml))
}
