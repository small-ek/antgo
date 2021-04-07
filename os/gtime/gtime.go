package gtime

import (
	"log"
	"time"
)

//时间类型转换

//Timestamp Get timestamp.时间转时间戳
func Timestamp(i time.Time) float64 {
	return float64(i.Unix())
}

//TimeString Time conversion String. 时间转string
func TimeString(i time.Time) string {
	if !i.IsZero() {
		return time.Unix(i.Unix(), 0).Format("2006-01-02 15:04:05")
	}
	return ""
}

//TimeFormat Time string conversion time. 字符串转时间
func Format(i interface{}, str string) time.Time {
	result, err := time.Parse(str, i.(string))
	if err != nil {
		log.Println(err.Error())
	}
	return result
}

//GetBeforeTime 获取之前的时间
func GetBefore(minusTimes string) time.Time {
	curTime := time.Now()
	dh, _ := time.ParseDuration("-" + minusTimes)
	timeStr := curTime.Add(dh)
	return timeStr
}

//GetAfterTime 获取之后的时间
func GetAfter(addTimes string) time.Time {
	curTime := time.Now()
	dh, _ := time.ParseDuration("+" + addTimes)
	timeStr := curTime.Add(dh)
	return timeStr
}
