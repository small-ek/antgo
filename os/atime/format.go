package atime

import (
	"errors"
	"fmt"
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

type (
	Date struct {
		time.Time
	}
)

//Now
func Now() Date {
	timeNow := time.Now()
	return WithTime(timeNow)
}

//WithTime
func WithTime(t time.Time) Date {
	return Date{t}
}

//WithTimestamp
func WithTimestamp(timestamp int64) Date {
	t := time.Unix(timestamp, 0)
	return WithTime(t)
}

//WithMillisecond ...
func WithMillisecond(millisecond int64) Date {
	sec := millisecond / 1e3
	nsec := millisecond % (millisecond / 1e3)
	t := time.Unix(sec, nsec)
	return WithTime(t)
}

//WithDate ...
func WithDate(year, month, date, hour, minute, second int) Date {
	t := time.Date(year, time.Month(month), date, hour, minute, second, 0, time.Local)
	return WithTime(t)
}

//Format ...
func (date Date) Format(layout string, chinese ...bool) (string, error) {
	var c bool
	if len(chinese) > 0 {
		c = chinese[0]
	}
	return date.parseLayout(layout, c)
}

//parseLayout [yyyy-MM-dd]->{{.year}}-{{.month}}-{{.day}}
func (date Date) parseLayout(layout string, chinese bool) (string, error) {
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
