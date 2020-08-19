package conv

import (
	"log"
	"time"
)

//时间类型转换

//获取时间戳
func GetTimestamp() float64 {
	return float64(time.Now().Unix())
}

//将<i>转换为sting
func TimeString(i time.Time) string {
	if !i.IsZero() {
		return time.Unix(i.Unix(), 0).Format("2006-01-02 15:04:05")
	}
	return ""
}

//将<i>转换为time.Time
func Time(i interface{}) time.Time {
	time, err := time.Parse("2006-01-02 15:04:05", i.(string))
	if err != nil {
		log.Println("时间转字符串失败" + err.Error())
	}
	return time
}
