package ant

import (
	"github.com/small-ek/antgo/db/mgo"
)

func Mgo(databaseName ...string) error {
	if err := mgo.InitEngine(); err != nil {
		return err
	}
	return nil
}
