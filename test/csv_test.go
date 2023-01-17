package test

import (
	"github.com/small-ek/antgo/os/csv"
	"github.com/small-ek/antgo/utils/conv"
	"log"
	"testing"
)

func TestCsv(t *testing.T) {
	acsv := csv.New("test.csv")
	for i := 0; i < 300; i++ {
		acsv.InsertOne([]string{conv.String(i) + "_12", "刘备_" + conv.String(i), "111", "Hello", "张飞", "关羽"})
	}
	acsv.Flush()
	log.Println(acsv.Read())
}
