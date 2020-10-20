package conv

import (
	"log"
	"time"
)

//时间类型转换

//GetTimestamp Get timestamp.
func GetTimestamp() float64 {
	return float64(time.Now().Unix())
}

//TimeString Time conversion String.
func TimeString(i time.Time) string {
	if !i.IsZero() {
		return time.Unix(i.Unix(), 0).Format("2006-01-02 15:04:05")
	}
	return ""
}

//Time Time string conversion time.
func Time(i interface{}) time.Time {
	result, err := time.Parse("2006-01-02 15:04:05", i.(string))
	if err != nil {
		log.Println(err.Error())
	}
	return result
}
