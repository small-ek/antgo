package test

import (
	"flag"
	"github.com/small-ek/antgo/validator"
	"log"
	"testing"
)

var Rule = map[string][]string{
	"password":     {"number:2|请输入大于8"},
	"old_password": {"require|请输入原有密码"},
}

type Tests struct {
	Name  string `json:"name" validate:"['require|请输入name','require|请输入name']"`
	Name2 string `json:"name2" validate:"require"`
}

func TestValidate(t *testing.T) {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	flag.Parse()
	var Request = map[string]interface{}{"old_password": "123456", "password": "12345.1212"}
	var err = validator.Default([]string{"old_password", "password"}, Rule).Check(Request)
	log.Print(err)

	var msg = map[string]string{
		"name.require":  "请输入name",
		"name.in":       "请输入正确的名称",
		"name2.require": "请输入name2",
		"name2.in":      "不在范围",
	}
	var err2 = validator.StructMsg(msg).CheckStruct(Tests{Name: "asaas", Name2: "2"})
	log.Println(err2)

}
