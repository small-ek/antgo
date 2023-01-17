package test

import (
	"github.com/small-ek/antgo/os/csv"
	"log"
	"testing"
)

func TestCsv(t *testing.T) {
	csvs := csv.New("test.csv")
	csvs.Insert([]string{"111", "刘备11", "23"})
	log.Println("------------------")
	log.Println(csvs.Read().Get())

}
