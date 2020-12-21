package validator

import (
	"errors"
	"github.com/small-ek/antgo/conv"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//New Validator default structure
type New struct {
	Rule    map[string][]string //Validation rules {"require|required", "max:25|maximum length","min:5|minimum length","number|number","email|mailbox",">:8|greater than 8"," <:8|Less than 8","=:8|equal to 8"}
	Scene   []string            //Detected field
	Message map[string]string
}

//Default
func Default(Scene []string, Rule map[string][]string) *New {
	return &New{
		Rule:  Rule,
		Scene: Scene,
	}
}

//SetRule
func (get *New) SetRule(Rule map[string][]string) *New {
	get.Rule = Rule
	return get
}

//SetScene
func (get *New) SetScene(Scene []string) *New {
	get.Scene = Scene
	return get
}

//CheckRule Form validator rules
func (get *New) Check(Request interface{}) error {
	var Scene = get.Scene
	var request map[string]interface{}
	conv.Struct(&request, Request)
	//循环要验证的数据
	for a := 0; a < len(Scene); a++ {
		var rowRule = get.Rule[Scene[a]]
		//循环需要验证的规则
		for b := 0; b < len(rowRule); b++ {
			if err := CheckItem(request[Scene[a]], rowRule[b]); err != nil {
				return err
			}
		}
	}
	return nil
}

//StructMessage
func StructMsg(msg map[string]string) *New {
	return &New{
		Message: msg,
	}
}

//CheckStruct TODO
func (get *New) CheckStruct(structModel interface{}) error {
	var types = reflect.TypeOf(structModel)
	if types.Kind() == reflect.Ptr {
		types = types.Elem()
	}
	var value = reflect.ValueOf(structModel)
	for i := 0; i < value.NumField(); i++ {
		var validate = types.Field(i).Tag.Get("validate")
		log.Print(validate)
		var jsonName = types.Field(i).Tag.Get("json")

		var validateStr = strings.Split(validate, "|")
		var val = value.Field(i).String()

		for j := 0; j < len(validateStr); j++ {
			var ruleType = validateStr[j]
			if strings.Index(ruleType, ":") > -1 {
				ruleType = strings.Split(ruleType, ":")[0]
			}
			var message = get.Message[jsonName+"."+ruleType]
			var err = CheckItem(val, validateStr[j]+"|"+message)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//CheckItem The validation field can define the rules you need TODO
func CheckItem(value interface{}, Rule string) error {
	var ruleSplit = strings.Split(Rule, "|")
	var rulesSplitStr = strings.Split(ruleSplit[0], ":")
	var rules = rulesSplitStr[0]
	var message = ruleSplit[1]
	var values = value.(string)
	if len(ruleSplit) == 1 {
		log.Println("Validator parameter syntax error")
		return nil
	}

	switch rules {
	//必填
	case "require":
		if isRequire(values) == false {
			return errors.New(message)
		}
	//在多少范围
	case "between":
		if between(values, rulesSplitStr[1]) == false {
			return errors.New(message)
		}
	//不在多少范围
	case "notBetween":
		if notBetween(values, rulesSplitStr[1]) == false {
			return errors.New(message)
		}
	//多少范围
	case "length":
		if length(values, rulesSplitStr[1]) == false {
			return errors.New(message)
		}
	//最大长度
	case "max":
		if isMax(values, rulesSplitStr[1]) == false {
			return errors.New(message)
		}
	//最小长度
	case "min":
		if isMin(values, rulesSplitStr[1]) == false {
			return errors.New(message)
		}
	//是否包含
	case "in":
		if in(values, rulesSplitStr[1]) == false {
			return errors.New(message)
		}
	//是否不包含
	case "notIn":
		if notIn(values, rulesSplitStr[1]) == false {
			return errors.New(message)
		}
	//数字类型包含小数点
	case "number":
		if isNumber(values) == false {
			return errors.New(message)
		}
	//邮箱
	case "email":
		if isEmail(values) == false {
			return errors.New(message)
		}
	//日期
	case "date":
		if date(values) == false {
			return errors.New(message)
		}
	//url验证
	case "url":
		if isUrl(values) == false {
			return errors.New(message)
		}
	//大于多少
	case ">":
		if moreThan(values, rulesSplitStr[1]) == false {
			return errors.New(message)
		}
	//小于多少
	case "<":
		if lessThan(values, rulesSplitStr[1]) == false {
			return errors.New(message)
		}
	//等于多少
	case "=":
		if equal(values, rulesSplitStr[1]) == false {
			return errors.New(message)
		}
	}

	return nil
}

//moreThan 大于
func moreThan(value, ruleStr string) bool {
	newValue, err := strconv.Atoi(value)

	if err != nil {
		log.Println(err.Error())
	}

	newValue2, err2 := strconv.Atoi(ruleStr)
	if err2 != nil {
		log.Println(err2.Error())
	}

	if newValue > newValue2 {
		return true
	}
	return false
}

//equal 等于
func equal(value, ruleStr string) bool {
	newValue, err := strconv.Atoi(value)

	if err != nil {
		log.Println(err.Error())
	}

	newValue2, err2 := strconv.Atoi(ruleStr)
	if err2 != nil {
		log.Println(err2.Error())
	}

	if newValue == newValue2 {
		return true
	}
	return false
}

//lessThan 小于
func lessThan(value, ruleStr string) bool {
	newValue, err := strconv.Atoi(ruleStr)

	if err != nil {
		log.Println(err.Error())
	}

	newValue2, err2 := strconv.Atoi(ruleStr)
	if err2 != nil {
		log.Println(err2.Error())
	}

	if newValue < newValue2 {
		return true
	}
	return false
}

//isRequire 验证数据不为空
func isRequire(value string) bool {
	if len(value) == 0 {
		return false
	}
	return true
}

//length 验证数据范围
func length(value, ruleStr string) bool {
	var str = strings.Split(ruleStr, ",")
	if len(str) == 1 && len(value) > conv.Int(str[0]) {
		return true
	}
	if len(str) == 2 && len(value) > conv.Int(str[0]) && len(value) < conv.Int(str[1]) {
		return true
	}
	return false
}

//isMax 验证数据最大长度
func isMax(value, ruleStr string) bool {
	newMax, err := strconv.Atoi(ruleStr)
	if err != nil {
		log.Println(err.Error())
	}
	if len(value) > newMax {
		return false
	}
	return true
}

//isMin 验证数据最小长度
func isMin(value, ruleStr string) bool {
	newMin, err := strconv.Atoi(ruleStr)

	if err != nil {
		log.Println(err.Error())
	}

	if len(value) < newMin {
		return false
	}
	return true
}

//isEmail 验证邮箱
func isEmail(value string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(value)
}

//isNumber 验证是否为合法数字
func isNumber(value string) bool {
	_, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return true
	}
	return false
}

//between 在多少之间
func between(value, ruleStr string) bool {
	var str = strings.Split(ruleStr, ",")
	if len(str) == 1 {
		return false
	}
	values, err := strconv.Atoi(value)

	if err != nil {
		log.Println(err.Error())
	}

	mins, err2 := strconv.Atoi(str[0])
	if err2 != nil {
		log.Println(err2.Error())
	}

	maxs, err3 := strconv.Atoi(str[1])
	if err2 != nil {
		log.Println(err3.Error())
	}

	if values > mins && values < maxs {
		return true
	}
	return false
}

//notBetween 不在多少之间
func notBetween(value, ruleStr string) bool {
	var str = strings.Split(ruleStr, ",")
	if len(str) == 1 {
		return false
	}
	values, err := strconv.Atoi(value)

	if err != nil {
		log.Println(err.Error())
	}

	mins, err2 := strconv.Atoi(str[0])
	if err2 != nil {
		log.Println(err2.Error())
	}

	maxs, err3 := strconv.Atoi(str[1])
	if err2 != nil {
		log.Println(err3.Error())
	}

	if values < mins && values > maxs {
		return true
	}
	return false
}

//date
func date(value string) bool {
	_, err := time.Parse("2006-01-02", value)
	if err != nil {
		return false
	}
	return true
}

//isUrl
func isUrl(value string) bool {
	if strings.Contains(value, "http") || strings.Contains(value, "https") || strings.Contains(value, "www") {
		return true
	}
	return false
}

//in
func in(value, ruleStr string) bool {
	if strings.Index(ruleStr, value) > -1 {
		return true
	}
	return false
}

//notIn
func notIn(value, ruleStr string) bool {
	if strings.Index(ruleStr, value) == -1 {
		return true
	}
	return false
}
