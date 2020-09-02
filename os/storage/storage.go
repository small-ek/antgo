package storage

import (
	"github.com/dgraph-io/badger/v2"
	"github.com/small-ek/ginp/conv"
	"log"
)

/*var Db *badger.DB
var err error*/

/*type Storage interface {
	Set(k, v []byte, expireAt int64) error

}*/
type Storage struct {
	Db *badger.DB
}

//打开数据库连接
func Open(path string) *Storage {
	db, err := badger.Open(badger.DefaultOptions(path))

	if err != nil {
		log.Println(err.Error())
	}
	return &Storage{Db: db}
}

//Set the cache data
func (this *Storage) Set(key string, value interface{}) error {
	err := this.Db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), conv.Bytes(value))
		return err
	})
	return err
}

func (this *Storage) Get(key string) ([]byte, error) {
	var valNot, valCopy []byte
	_ = this.Db.View(func(txn *badger.Txn) error {
		item, _ := txn.Get([]byte(key))
		_ = item.Value(func(val []byte) error {
			// Copying or parsing val is valid.
			valCopy = append([]byte{}, val...)
			// Assigning val slice to another variable is NOT OK.
			valNot = val // Do not do this.
			return nil
		})

		valCopy, _ = item.ValueCopy(nil)
		return nil
	})
	log.Println(valNot)
	return valCopy, nil
}

/*
func Get(key string) {
	err := Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("answer"))
		handle(err)

		var valNot, valCopy []byte
		err := item.Value(func(val []byte) error {
			// This func with val would only be called if item.Value encounters no error.

			// Accessing val here is valid.
			fmt.Printf("The answer is: %s\n", val)

			// Copying or parsing val is valid.
			valCopy = append([]byte{}, val...)

			// Assigning val slice to another variable is NOT OK.
			valNot = val // Do not do this.
			return nil
		})


		// DO NOT access val here. It is the most common cause of bugs.
		fmt.Printf("NEVER do this. %s\n", valNot)

		// You must copy it to use it outside item.Value(...).
		fmt.Printf("The answer is: %s\n", valCopy)

		// Alternatively, you could also use item.ValueCopy().
		valCopy, err = item.ValueCopy(nil)

		fmt.Printf("The answer is: %s\n", valCopy)

		return err
	})
}*/
