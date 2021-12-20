package db

import (
	"github.com/small-ek/antgo/conv"
	"github.com/small-ek/antgo/os/config"
	loggers "github.com/small-ek/antgo/os/logger"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var Master *gorm.DB
var *gorm.DBResolver
var datetimePrecision = 2

type Db struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	Params   string `json:"params"`
	Log      bool   `json:"log"`
}

func InitDb() {
	cfg := config.Decode()
	connections := cfg.Get("connections").Maps()
	default_connections := cfg.Get("system.default_connections").String()

	for i := 0; i < len(connections); i++ {
		value := connections[i]
		row := Db{}
		conv.Struct(&row, value)

		switch row.Type {
		case "mysql":
			var dns = row.Username + ":" + row.Password + "@tcp(" + row.Hostname + ":" + row.Port + ")/" + row.Database + "?" + row.Params

			if row.Name == default_connections {
				row.Open(Mysql(dns), getConfig(row.Log))
			} else {
				dbresolver.Register(dbresolver.Config{
					Replicas: []gorm.Dialector{mysql.Open(dns)},
					// sources/replicas 负载均衡策略
					Policy: dbresolver.RandomPolicy{},
				}, row.Name)
			}

			break
		case "pgsql":
			var dns = row.Username + ":" + row.Password + "@tcp(" + row.Hostname + ":" + row.Port + ")/" + row.Database + "?" + row.Params

			if row.Name == default_connections {
				row.Open(Postgres(dns), getConfig(row.Log))
			} else {
				row.Use(dbresolver.Register(dbresolver.Config{
					Replicas: []gorm.Dialector{mysql.Open(dns)},
					// sources/replicas 负载均衡策略
					Policy: dbresolver.RandomPolicy{},
				}, row.Name))
			}

			break
		}
	}

}

//getConfig
func getConfig(isLog bool) *gorm.Config {
	if isLog {
		loggers := New(zap.L())
		loggers.SetAsDefault()
		return &gorm.Config{
			Logger:                                   loggers.LogMode(4),
			DisableForeignKeyConstraintWhenMigrating: true,
		}
	} else {
		return &gorm.Config{}
	}
}

//Open connection
func (d *Db) Open(Dialector gorm.Dialector, opts gorm.Option) {
	var db, err = gorm.Open(Dialector, opts)
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
	return postgres.New(postgres.Config{
		DriverName:           "",
		DSN:                  dns,
		PreferSimpleProtocol: false,
		WithoutReturning:     false,
		Conn:                 nil,
	})
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
