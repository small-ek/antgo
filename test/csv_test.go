package test

import (
	"github.com/small-ek/antgo/os/acsv"
	"github.com/small-ek/antgo/utils/conv"
	"testing"
)

func TestCsv(t *testing.T) {
	c, _ := acsv.New("test.csv")
	c2, _ := acsv.New("test2.csv")

	for i := 0; i < 40000; i++ {
		c.AddRow([]string{conv.String(i), "John Doe", "john.doe@example.com"})
		c2.AddRow([]string{conv.String(i), "John Doe2", "john.doe@example.com2"})
	}
	c.Write()
	c2.Write()
}
