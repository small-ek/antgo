package test

import (
	"github.com/small-ek/antgo/os/atime"
	"log"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	// 通过时间字符串创建
	log.Println(atime.New("2020-10-24 12:00:00"))
	// 通过time.Time对象创建
	log.Println(atime.New(time.Now()))
	// 通过时间戳(秒)创建
	log.Println(atime.New(1603710586))
	// 通过时间戳(纳秒)创建
	log.Println(atime.New(1603710586660409000))
	d := atime.Now()

	log.Println(d.String())
	log.Println(d.Adds(-time.Hour * 2).Format("yyyy-MM-dd HH:mm:ss"))
	log.Println(d.Timestamp())
	log.Println(d.Millisecond())
	log.Println(d.AddDates(0, 1, 0).Format("yyyy-MM-dd HH:mm:ss"))
	format := d.Format("yyyy-MM-dd HH:mm:ss")
	log.Println(format)
}
