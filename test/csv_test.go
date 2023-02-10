package test

import (
	"github.com/small-ek/antgo/os/csv"
	"github.com/small-ek/antgo/utils/conv"
	"testing"
)

func TestCsv(t *testing.T) {
	c, _ := csv.NewCSV("test.csv")
	c2, _ := csv.NewCSV("test2.csv")

	for i := 0; i < 40000000; i++ {
		c.AddRow([]string{conv.String(i), "John Doe", "john.doe@example.com"})
		c2.AddRow([]string{conv.String(i), "John Doe2", "john.doe@example.com2"})
	}
	c.Write()
	c2.Write()
}
