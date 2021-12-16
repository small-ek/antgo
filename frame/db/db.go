package db

import (
	"github.com/small-ek/antgo/os/config"
	loggers "github.com/small-ek/antgo/os/logger"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"log"
)

var Master *gorm.DB
var datetimePrecision = 2

type Db struct {
	Name     string
	Type     string
	hostname string
	port     string
	username string
	password string
	database string
	params   string
	log      bool
}

func InitDb() {
	cfg := config.Decode()
	connections := cfg.Get("connections").Maps()
	log.Println(connections)
	for i := 0; i < len(connections); i++ {
		row := connections[i]
		log.Println(row)
		switch row["type"] {
		case "mysql":
			var dns = row["username"].(string) + ":" + row["password"].(string) + "@tcp(" + row["hostname"].(string) + ":" + row["port"].(string) + ")/" + row["database"].(string) + "?" + row["params"].(string)
			if 0 == i {
				db, err := gorm.Open(Mysql(dns))
				if err != nil {
					loggers.Write.Panic(err.Error())
				}
				Master = db
			}

			break
		case "pgsql":
			var dns = row["username"].(string) + ":" + row["password"].(string) + "@tcp(" + row["hostname"].(string) + ":" + row["port"].(string) + ")/" + row["database"].(string) + "?" + row["params"].(string)
			if 0 == i {
				db, err := gorm.Open(Postgres(dns))
				if err != nil {
					loggers.Write.Panic(err.Error())
				}

				Master = db
			}
			break
		}
	}
}

//getConfig
func getConfig(log bool) *gorm.Config {
	logger := New(zap.L())
	logger.SetAsDefault()
	if log {
		return &gorm.Config{
			Logger: logger,
		}
	} else {
		return &gorm.Config{
			Logger: logger,
		}
	}
}

//Open connection
func (d *Db) Open(Dialector gorm.Dialector) {
	var db, err = gorm.Open(Dialector)
	if err != nil {
		loggers.Write.Panic(err.Error())
	}
	Master = db
}

//Use
func (d *Db) Use(plugin gorm.Plugin) {
	if err := Master.Use(plugin); err != nil {
		loggers.Write.Panic(err.Error())
	}

}

//Mysql connection
func Mysql(dns string) gorm.Dialector {
	return mysql.New(mysql.Config{
		DSN:                       dns, // DSN data source name
		DefaultStringSize:         256, // string 类型字段的默认长度
		DefaultDatetimePrecision:  &datetimePrecision,
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	})
}

//Postgres connection
func Postgres(dns string) gorm.Dialector {
	return postgres.Open(dns)
}

//Distributed
func (d *Db) Distributed(config dbresolver.Config, datas ...interface{}) *dbresolver.DBResolver {
	return dbresolver.Register(config, datas)
}

//Close 关闭数据库
func Close() {
	var db, err = Master.DB()
	if err != nil {
		loggers.Write.Panic(err.Error())
	}
	db.Close()
}
