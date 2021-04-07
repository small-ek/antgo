package test

import (
	"github.com/small-ek/antgo/os/gtime"
	"log"
	"testing"
)

func TestTime(t *testing.T) {
	log.Println(gtime.GetBefore("24h").Format("2006-01-02"))
	log.Println(gtime.GetAfter("24h").Format("2006-01-02"))
	log.Println(gtime.Format("2022-01-02", "2006-01-02"))
}
