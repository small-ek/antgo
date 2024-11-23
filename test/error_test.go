package test

import (
	"database/sql"
	"github.com/small-ek/antgo/os/aerror"
	"log"
	"testing"
)

func foo() error {
	return aerror.Wrap(sql.ErrNoRows, "foo failed")
}

func bar() error {
	err := foo()
	return aerror.Wrap(err, "bar failed")
}

func baz() error {
	err := bar()
	return aerror.Wrap(err, "baz failed")
}
func TestError(t *testing.T) {
	err := baz()
	log.Println(err.Error())
}
