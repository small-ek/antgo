package test

import (
	"github.com/small-ek/antgo/os/csv"
	"github.com/small-ek/antgo/utils/conv"
	"log"
	"testing"
)

func TestCsv(t *testing.T) {

	for i := 0; i < 90000; i++ {
		err := csv.New("test.csv").InsertOne([]string{conv.String(i) + "_12", "刘备_" + conv.String(i), "111", "Hello", "张飞", "关羽"})
		if err != nil {
			return
		}
	}
	log.Println(csv.New("test.csv").GetCount())
	log.Println(csv.New("test.csv").Read())

}
