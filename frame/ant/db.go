package ant

import (
	"github.com/small-ek/antgo/frame/db"
	"gorm.io/gorm"
)

// Db Get database connection
func Db() *gorm.DB {
	return db.Master
}
