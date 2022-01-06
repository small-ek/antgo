package test

import (
	"flag"
	"github.com/small-ek/antgo/utils/validator"
	"log"
	"testing"
)

var Rule = map[string][]string{
	"password":     {"number:2|请输入大于8"},
	"old_password": {"require|请输入原有密码"},
}

type Tests struct {
	Name  int    `json:"name" validate:"require|请输入name;number:2|请输入name2"`
	Name2 string `json:"name2" validate:"require|请输入name2"`
}

func TestValidate(t *testing.T) {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	flag.Parse()

	var Request = map[string]interface{}{"old_password": "123456", "password": "12345.1212"}
	var err = validator.New([]string{"old_password", "password"}, Rule).Check(Request)
	log.Println(err)

	var err2 = validator.CheckStruct(Tests{Name: 11111, Name2: "2222"})
	log.Println(err2)
}
