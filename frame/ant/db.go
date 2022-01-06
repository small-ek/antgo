package ant

import (
	"github.com/small-ek/antgo/db/adb"
	"gorm.io/gorm"
)

// Db Get database connection
func Db() *gorm.DB {
	return adb.Master
}

// CloseDb Close database connection
func CloseDb() {
	adb.Close()
}
