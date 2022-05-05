package adb

import (
	"github.com/small-ek/antgo/os/config"
	loggers "github.com/small-ek/antgo/os/logger"
	"github.com/small-ek/antgo/utils/conv"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"log"
)

var Master *gorm.DB
var Resolver *dbresolver.DBResolver
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
	Dns      string `json:"dns"`
}

func InitDb() {
	cfg := config.Decode()
	connections := cfg.Get("connections").Maps()
	default_connections := cfg.Get("system.default_connections").String()

	if default_connections == "" {
		log.Println("default_connections is empty")
		return
	}

	for i := 0; i < len(connections); i++ {
		value := connections[i]
		row := Db{}
		conv.Struct(&row, value)

		switch row.Type {
		case "mysql":
			dns := row.Username + ":" + row.Password + "@tcp(" + row.Hostname + ":" + row.Port + ")/" + row.Database + "?" + row.Params
			dbTypeSelect(row, default_connections, Mysql(dns))
			break

		case "postgres":
			dns := "host=" + row.Hostname + " port=" + row.Port + " user=" + row.Username + " dbname=" + row.Database + " " + row.Params + " password=" + row.Password
			dbTypeSelect(row, default_connections, Postgres(dns))
			break

		case "mssql":
			dns := "sqlserver://" + row.Username + ":" + row.Password + "@" + row.Hostname + ":" + row.Port + "?database=" + row.Database
			dbTypeSelect(row, default_connections, Sqlserver(dns))
			break

		case "sqlite":
			dns := row.Dns
			dbTypeSelect(row, default_connections, Sqlite(dns))
			break
		}
	}

	if len(connections) == 0 {
		return
	}

	err := Master.Use(Resolver)
	if err != nil {
		panic(err)
	}
}

func dbTypeSelect(row Db, default_conns string, dialector gorm.Dialector) {
	if row.Name == default_conns {
		row.Open(dialector, getConfig(row.Log))
	} else {
		Resolver = dbresolver.Register(dbresolver.Config{
			Replicas: []gorm.Dialector{dialector},
			// sources/replicas 负载均衡策略
			Policy: dbresolver.RandomPolicy{},
		}, row.Name)
	}
}

//getConfig
func getConfig(isLog bool) *gorm.Config {
	if isLog {
		zapLog := New(zap.L())
		zapLog.SetAsDefault()
		return &gorm.Config{
			Logger:                                   zapLog.LogMode(4),
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
		PreferSimpleProtocol: true,
		WithoutReturning:     false,
		Conn:                 nil,
	})
}

//Sqlserver connection
func Sqlserver(dns string) gorm.Dialector {
	return sqlserver.Open(dns)
}

//Sqlite connection
func Sqlite(dns string) gorm.Dialector {
	return sqlite.Open(dns)
}

//Distributed
func (d *Db) Distributed(config dbresolver.Config, datas ...interface{}) *dbresolver.DBResolver {
	return dbresolver.Register(config, datas)
}

//Close 关闭数据库
func Close() {
	var db, err = Master.DB()
	if err != nil {
		loggers.Write.Error(err.Error())
	}

	if db != nil {
		if err2 := db.Close(); err2 != nil {
			return
		}
	}
}
