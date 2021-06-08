package atime

import (
	"errors"
	"fmt"
	"github.com/small-ek/antgo/conv"
	"github.com/small-ek/antgo/os/logs"
	"log"
	"strings"
	"time"
)

var (
	ErrLayout = errors.New(`parse layout failed`)
	monthAbbr = [...]string{
		"Jan",
		"Feb",
		"Mar",
		"Apr",
		"May",
		"Jun",
		"Jul",
		"Aug",
		"Sept",
		"Oct",
		"Nov",
		"Dec",
	}
	weekDayAbbr = [...]string{
		"Sun",
		"Mon",
		"Sun",
		"Tue",
		"Wed",
		"Thur",
		"Fri",
		"Sat",
	}
	weekDayChinese = [...]string{
		"星期日",
		"星期一",
		"星期二",
		"星期三",
		"星期四",
		"星期五",
		"星期六",
	}
)

type Times struct {
	time.Time
}

//New 创建对象
func New(param ...interface{}) *Times {
	if len(param) > 0 {
		switch r := param[0].(type) {
		case time.Time:
			return WithTime(r)
		case *time.Time:
			return WithTime(*r)
		case string:
			return StrToTime(r)
		case []byte:
			return StrToTime(string(r))
		case int:
			return NewFromTimeStamp(int64(r))
		case int64:
			return NewFromTimeStamp(r)
		default:
			return NewFromTimeStamp(conv.Int64(r))
		}
	}
	return &Times{time.Now()}
}

//StrToTime String转Time
func StrToTime(str string) *Times {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", str, time.Local)
	if err != nil {
		log.Println(err)
	}
	return &Times{t}
}

// NewFromTimeStamp creates and returns a Time object with given timestamp,
// which can be in seconds to nanoseconds.
// Eg: 1600443866 and 1600443866199266000 are both considered as valid timestamp number.
func NewFromTimeStamp(timestamp int64) *Times {
	if timestamp == 0 {
		return &Times{}
	}
	var sec, nano int64
	if timestamp > 1e9 {
		for timestamp < 1e18 {
			timestamp *= 10
		}
		sec = timestamp / 1e9
		nano = timestamp % 1e9
	} else {
		sec = timestamp
	}
	return WithTime(time.Unix(sec, nano))
}

//Now
func Now() *Times {
	timeNow := time.Now()
	return WithTime(timeNow)
}

//WithTime
func WithTime(t time.Time) *Times {
	return &Times{t}
}

//Format ...
func (t *Times) Format(layout string, chinese ...bool) string {
	var c bool
	if len(chinese) > 0 {
		c = chinese[0]
	}
	d, err := t.parseLayout(layout, c)
	if err != nil {
		logs.Error(err.Error())
	}
	return d
}

//parseLayout [yyyy-MM-dd]->{{.year}}-{{.month}}-{{.day}}
func (date *Times) parseLayout(layout string, chinese bool) (string, error) {
	if len(strings.TrimSpace(layout)) == 0 {
		return "", ErrLayout
	}
	ti := time.Unix(date.Unix(), date.UnixNano()/1e6%date.Unix())
	year, monthNumber, dayOfMonth := date.Date()
	thisMonthFirstDay := time.Date(year, monthNumber, 1, 0, 0, 0, 0, time.Local)
	thisYearFirstDay := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	yearRight := year % 100
	monthCom := monthNumber.String()
	monthAbbr := monthAbbr[monthNumber-1]
	weekOfMonth := (date.Day()-1)/7 + 1
	weekdayOfThisMonthFirstDay := thisMonthFirstDay.Weekday()
	relativeWeekOfMonth := (dayOfMonth+int(weekdayOfThisMonthFirstDay-1))/7 + 1

	dayOfYear := date.YearDay()
	dayOfWeek := date.Weekday()
	weekOfYear := (dayOfYear+int(thisYearFirstDay.Weekday()-1))/7 + 1
	if weekOfYear/53 >= 1 {
		weekOfYear = 1
	}
	hour, minute, second := date.Clock()
	unix := date.Unix()
	var millisecond int64
	if unix > 0 {
		unixNano := date.UnixNano()
		millisecond = unixNano % unix
	}
	rfc822z := ti.Format("-0700")
	stz := ti.Format("MST")
	var am bool
	am = date.Hour() < 12

	var t = new(strings.Builder)
	var i = 0

	for i < len(layout) {
		c := layout[i]
		switch c {
		case 'y': // 年[year]
			y, endIndex := end(i, layout, 'y')
			if length := len(y); length > 3 {
				t.WriteString(fmt.Sprintf("%d", year))
			} else {
				t.WriteString(fmt.Sprintf("%0*d", 2, yearRight))
			}
			i = endIndex
		case 'M': // 月[month]
			m, endIndex := end(i, layout, 'M')
			if length := len(m); length > 3 {
				t.WriteString(monthCom)
			} else if length == 3 {
				t.WriteString(monthAbbr)
			} else {
				validLength := int(monthNumber / 10)
				t.WriteString(fmt.Sprintf("%0*d", 2-validLength, monthNumber))
			}
			i = endIndex
		case 'w': // 年中的周数[number]
			w, endIndex := end(i, layout, 'w')
			validLength := len(w)
			t.WriteString(fmt.Sprintf("%0*d", validLength, weekOfYear))
			i = endIndex
		case 'W': // 月份中的周数[number]
			W, endIndex := end(i, layout, 'W')
			validLength := len(W)
			t.WriteString(fmt.Sprintf("%0*d", validLength, relativeWeekOfMonth))
			i = endIndex
		case 'D': // 年中的天数[number]
			D, endIndex := end(i, layout, 'D')
			validLength := len(D)
			t.WriteString(fmt.Sprintf("%0*d", validLength, dayOfYear))
			i = endIndex
		case 'd': // 月份中的天数[number]
			d, endIndex := end(i, layout, 'd')
			validLength := len(d)
			t.WriteString(fmt.Sprintf("%0*d", validLength, dayOfMonth))
			i = endIndex
		case 'F': // 月份中的星期[number]
			F, endIndex := end(i, layout, 'F')
			var numberLength = 1
			if weekOfMonth > 9 {
				numberLength = 2
			}
			validLength := len(F) - numberLength + 1
			t.WriteString(fmt.Sprintf("%0*d", validLength, weekOfMonth))
			i = endIndex
		case 'E': // 星期中的天数[text]
			E, endIndex := end(i, layout, 'E')
			if chinese {
				t.WriteString(weekDayChinese[dayOfWeek])
			} else {
				if length := len(E); length > 3 {
					t.WriteString(dayOfWeek.String())
				} else if length == 3 {
					t.WriteString(weekDayAbbr[dayOfWeek])
				} else {
					t.WriteString(fmt.Sprintf("%0*d", length, dayOfWeek))
				}
			}
			i = endIndex
		case 'a': // am/pm[text]
			_, endIndex := end(i, layout, 'a')
			if chinese {
				if am {
					t.WriteString("上午")
				} else {
					t.WriteString("下午")
				}
			} else {
				if am {
					t.WriteString("AM")
				} else {
					t.WriteString("PM")
				}
			}
			i = endIndex
		case 'H': // 一天中的小时数，0-23[number]
			H, endIndex := end(i, layout, 'H')
			validLength := len(H)
			t.WriteString(fmt.Sprintf("%0*d", validLength, hour))
			i = endIndex
		case 'k': // 一天中的小时数，1-24[number]
			k, endIndex := end(i, layout, 'k')
			validLength := len(k)
			if hour == 0 {
				t.WriteString(fmt.Sprintf("%0*d", validLength, 1))
			} else {
				t.WriteString(fmt.Sprintf("%0*d", validLength, hour))
			}
			i = endIndex
		case 'K': // am/pm小时数，0-11[number]
			K, endIndex := end(i, layout, 'K')
			validLength := len(K)
			t.WriteString(fmt.Sprintf("%0*d", validLength, hour%12))
			i = endIndex
		case 'h': // am/pm小时数,1-12[number]
			h, endIndex := end(i, layout, 'h')
			validLength := len(h)
			if hour == 0 {
				t.WriteString(fmt.Sprintf("%0*d", validLength, 1))
			} else {
				t.WriteString(fmt.Sprintf("%0*d", validLength, (hour)%12))
			}
			i = endIndex
		case 'm': // 小时中的分钟数[number]
			m, endIndex := end(i, layout, 'm')
			validLength := len(m)
			t.WriteString(fmt.Sprintf("%0*d", validLength, minute))
			i = endIndex
		case 's': // 分钟中的秒数[number]
			s, endIndex := end(i, layout, 's')
			validLength := len(s)
			t.WriteString(fmt.Sprintf("%0*d", validLength, second))
			i = endIndex
		case 'S': // 毫秒数[number]
			S, endIndex := end(i, layout, 'S')
			validLength := len(S)
			t.WriteString(fmt.Sprintf("%0*d", validLength, millisecond))
			i = endIndex
		case 'z': // 时区（General）
			_, endIndex := end(i, layout, 'z')
			t.WriteString(stz)
			i = endIndex
		case 'Z': // 时区（RFC）
			_, endIndex := end(i, layout, 'Z')
			t.WriteString(rfc822z)
			i = endIndex
		default:
			t.WriteByte(c)
			i = i + 1
		}
	}
	return t.String(), nil
}

//end ...
func end(from int, in string, target rune) (string, int) {
	var out = new(strings.Builder)
	for i := from; i < len(in); i++ {
		r := rune(in[i])
		from = i
		if r == target {
			out.WriteRune(r)
			if i == len(in)-1 {
				return out.String(), i + 1
			}
			continue
		}
		return out.String(), i
	}
	return "", from + 1
}
