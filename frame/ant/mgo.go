package ant

import "github.com/small-ek/antgo/db/mgo"

func Mgo(databaseName ...string) *mgo.Mgo {
	return mgo.Connect(databaseName...)
}
