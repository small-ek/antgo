package test

import (
	"github.com/small-ek/antgo/os/csv"
	"log"
	"testing"
)

func TestCsv(t *testing.T) {
	csvs := csv.New("test.csv")
	csvs.Create().Insert([][]string{
		{"1", "刘备11", "23"},
		{"2", "张飞", "23"},
		{"3", "关羽", "23"},
		{"4", "赵云", "23"},
		{"5", "黄忠", "23"},
		{"6", "马超", "23"},
	})
	log.Println("------------------")
	log.Println(csvs.Read().Get())

}
