package atime

import (
	"time"
)

var layout = "2006-01-02 15:04:05"

//Timestamp Get timestamp.时间转时间戳
func Timestamp(i time.Time) int64 {
	return i.Unix()
}

//Millisecond ...
func Millisecond(i time.Time) int64 {
	return i.UnixNano() / 1e6
}

//TimeString Time conversion String. 时间转string
func TimeString(i time.Time) string {
	if !i.IsZero() {
		return time.Unix(i.Unix(), 0).Format(layout)
	}
	return ""
}

//GetBefore 获取之前的时间
func GetBefore(minusTimes string) time.Time {
	curTime := time.Now()
	dh, _ := time.ParseDuration("-" + minusTimes)
	timeStr := curTime.Add(dh)
	return timeStr
}

//GetAfter 获取之后的时间
func GetAfter(addTimes string) time.Time {
	curTime := time.Now()
	dh, _ := time.ParseDuration("+" + addTimes)
	timeStr := curTime.Add(dh)
	return timeStr
}
