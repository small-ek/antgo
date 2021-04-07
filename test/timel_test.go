package test

import (
	"github.com/small-ek/antgo/os/atime"
	"log"
	"testing"
)

func TestTime(t *testing.T) {
	log.Println(atime.GetBefore("24h").Format("2006-01-02"))
	log.Println(atime.GetAfter("24h").Format("2006-01-02"))

	d := atime.Now()
	log.Println(atime.TimeString(d.Time))
	log.Println(atime.Millisecond(d.Time))
	log.Println(atime.Timestamp(d.Time))
	log.Println(d.YearDay())
	ret, err := d.Format("yyyy-MM-dd HH:mm:ss")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(ret)
}
