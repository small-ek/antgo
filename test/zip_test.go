package test

import (
	"github.com/small-ek/ginp/encoding/zip"
	"log"
	"testing"
)

func TestZip(t *testing.T) {
	// List of Files to Zip
	files := []string{"test.txt", "test2.txt"}
	output := "done.zip"

	if err := zip.Create(output, files); err != nil {
		panic(err)
	}

	files, err := zip.Unzip(output, "done")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Zipped File:", output)
}
