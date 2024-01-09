package ant

import (
	"github.com/small-ek/antgo/db/adb"
	"gorm.io/gorm"
)

// Db Get database connection
func Db(name string) *gorm.DB {
	return adb.Master[name]
}

// CloseDb Close database connection
func CloseDb() {
	adb.Close()
}
