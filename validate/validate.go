package validate

import (
	"errors"
	"github.com/small-ek/ginp/conv"
	"log"
	"regexp"
	"strconv"
	"strings"
)

//New Validator default structure
type New struct {
	Request interface{}         //Request data
	Rule    map[string][]string //Validation rules {"require|required", "max:25|maximum length","min:5|minimum length","number|number","email|mailbox",">:8|greater than 8"," <:8|Less than 8","=:8|equal to 8"}
	Scene   []string            //Detected field
}

//CheckRule Form validator rules
func (this *New) CheckRule() error {
	var Scene = this.Scene
	var request map[string]string
	conv.Struct(&request, this.Request)
	//循环要验证的数据
	for a := 0; a < len(Scene); a++ {
		var rowRule = this.Rule[Scene[a]]
		//循环需要验证的规则
		for b := 0; b < len(rowRule); b++ {
			if err := CheckItem(request[Scene[a]], rowRule[b]); err != nil {
				return err
			}
		}

	}

	return nil
}

//TODO
//CheckItem The validation field can define the rules you need
func CheckItem(value string, Rule string) error {
	var message = strings.Split(Rule, "|")
	var newValue = strings.Split(message[0], ":")
	if len(message) == 1 {
		log.Println("验证器参数不正确")
		return errors.New("验证器参数不正确")
	}

	switch newValue[0] {
	//必填
	case "require":
		if isRequire(value) == false {
			return errors.New(message[1])
		}
	//最大长度
	case "max":
		if isMax(value, newValue[1]) == false {
			return errors.New(message[1])
		}
	//最小长度
	case "min":
		if isMin(value, newValue[1]) == false {
			return errors.New(message[1])
		}
	//数字类型包含小数点
	case "number":
		if isNumber(value) == false {
			return errors.New(message[1])
		}
	//邮箱
	case "email":
		if isEmail(value) == false {
			return errors.New(message[1])
		}
	//大于多少
	case ">":
		if moreThan(value, newValue[1]) == false {
			return errors.New(message[1])
		}
	//小于多少
	case "<":
		if lessThan(value, newValue[1]) == false {
			return errors.New(message[1])
		}
	//等于多少
	case "=":
		if equal(value, newValue[1]) == false {
			return errors.New(message[1])
		}
	}

	return nil
}

//大于
func moreThan(value, value2 string) bool {
	newValue, err := strconv.Atoi(value)

	if err != nil {
		log.Println("数据类型不正确")
	}

	newValue2, err2 := strconv.Atoi(value2)
	if err2 != nil {
		log.Println("数据类型不正确")
	}

	if newValue > newValue2 {
		return true
	}
	return false
}

//等于
func equal(value, value2 string) bool {
	newValue, err := strconv.Atoi(value)

	if err != nil {
		log.Println("数据类型不正确")
	}

	newValue2, err2 := strconv.Atoi(value2)
	if err2 != nil {
		log.Println("数据类型不正确")
	}

	if newValue == newValue2 {
		return true
	}
	return false
}

//小于
func lessThan(value, value2 string) bool {
	newValue, err := strconv.Atoi(value)

	if err != nil {
		log.Println("数据类型不正确")
	}

	newValue2, err2 := strconv.Atoi(value2)
	if err2 != nil {
		log.Println("数据类型不正确")
	}

	if newValue < newValue2 {
		return true
	}
	return false
}

//验证数据不为空
func isRequire(value string) bool {
	if len(value) == 0 {
		return false
	}
	return true
}

//验证数据最大长度
func isMax(value, max string) bool {
	newMax, err := strconv.Atoi(max)
	if err != nil {
		log.Println("数据类型不正确")
	}
	if len(value) > newMax {
		return false
	}
	return true
}

//验证数据最小长度
func isMin(value, min string) bool {
	newMin, err := strconv.Atoi(min)

	if err != nil {
		log.Println("数据类型不正确")
	}

	if len(value) < newMin {
		return false
	}
	return true
}

//验证邮箱
func isEmail(value string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(value)
}

//验证是否为合法数字
func isNumber(val interface{}) bool {
	switch val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
	case float32, float64, complex64, complex128:
		return true
	case string:
		str := val.(string)
		if str == "" {
			return false
		}
		// Trim any whitespace
		str = strings.Trim(str, " \\t\\n\\r\\v\\f")
		if str[0] == '-' || str[0] == '+' {
			if len(str) == 1 {
				return false
			}
			str = str[1:]
		}
		// hex
		if len(str) > 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X') {
			for _, h := range str[2:] {
				if !((h >= '0' && h <= '9') || (h >= 'a' && h <= 'f') || (h >= 'A' && h <= 'F')) {
					return false
				}
			}
			return true
		}
		// 0-9,Point,Scientific
		p, s, l := 0, 0, len(str)
		for i, v := range str {
			if v == '.' { // Point
				if p > 0 || s > 0 || i+1 == l {
					return false
				}
				p = i
			} else if v == 'e' || v == 'E' { // Scientific
				if i == 0 || s > 0 || i+1 == l {
					return false
				}
				s = i
			} else if v < '0' || v > '9' {
				return false
			}
		}
		return true
	}

	return false
}
