package test

import (
	"errors"
	"log"
	"testing"
)

func TestErr(t *testing.T) {



	var err = errors.New("2222")
	if err != nil {
		log.Panicln(err)
	}
	err = errors.New("1111")
	if err != nil {
		log.Println(err)
	}
}
func a() error {
	return nil
}
