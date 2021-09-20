package test

import (
	"fmt"
	"github.com/small-ek/antgo/encoding/axml"
	"log"
	"testing"
)

var xmlStr string = `
<?xml version="1.0" encoding="UTF-8"?>
<resources>
	<string name="VideoLoading">Loading video…</string>
	<string2 name="ApplicationName">这是新的ApplicationName</string2>
</resources>
`

func TestXml(t *testing.T) {
	userMap := make(map[string]string)
	userMap["name"] = "Name"
	userMap["id"] = "1"

	buf, _ := axml.Encode(userMap)
	fmt.Println(string(buf))

	result, err := axml.Decode([]byte(xmlStr))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)

	result2, err := axml.ToJson([]byte(xmlStr))
	log.Println(string(result2))
	log.Println(err)
}
