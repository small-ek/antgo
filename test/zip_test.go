package test

import (
	"github.com/small-ek/antgo/encoding/azip"
	"log"
	"testing"
)

func TestZip(t *testing.T) {
	// List of Files to Zip
	files := []string{"test.txt", "test2.txt"}
	output := "done.zip"

	if err := azip.Create(output, files); err != nil {
		panic(err)
	}

	files2, err2 := azip.Unzip(output, "done")
	if err2 != nil {
		log.Println(err2)
	}
	log.Println(files2)
	log.Println("Zipped File:", output)
}
