package adb

import (
	"github.com/small-ek/antgo/os/alog"
	"github.com/small-ek/antgo/os/config"
	"github.com/small-ek/antgo/utils/conv"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"strconv"
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
	Dsn      string `json:"dsn"`
}

func InitDb() {

	connections := conv.Map(config.Get("connections"))
	default_connections := conv.String(config.Get("system.default_connections"))
	if default_connections != "" {

		for i := 0; i < len(connections); i++ {
			value := connections[strconv.Itoa(i)]
			row := Db{}
			conv.Struct(&row, value)
			dsn := row.Dsn

			switch row.Type {
			case "mysql":
				if row.Dsn == "" {
					dsn = row.Username + ":" + row.Password + "@tcp(" + row.Hostname + ":" + row.Port + ")/" + row.Database + "?" + row.Params
				}

				if row.Name == default_connections {
					row.Open(Mysql(dsn), getConfig(row.Log))
				} else {
					Resolver = dbresolver.Register(dbresolver.Config{
						Replicas: []gorm.Dialector{Mysql(dsn)},
						// sources/replicas 负载均衡策略
						Policy: dbresolver.RandomPolicy{},
					}, row.Name)
				}
				break
			case "pgsql":
				if row.Dsn == "" {
					dsn = "host=" + row.Hostname + " port=" + row.Port + " user=" + row.Username + " dbname=" + row.Database + " " + row.Params + " password=" + row.Password + row.Params
				}
				if row.Name == default_connections {
					row.Open(Postgres(dsn), getConfig(row.Log))
				} else {
					Resolver = dbresolver.Register(dbresolver.Config{
						Replicas: []gorm.Dialector{Postgres(dsn)},
						// sources/replicas 负载均衡策略
						Policy: dbresolver.RandomPolicy{},
					}, row.Name)
				}
				break
			case "sqlsrv":
				if row.Dsn == "" {
					dsn = "sqlserver://" + row.Username + ":" + row.Password + "@" + row.Hostname + ":" + row.Port + "?database=" + row.Database + row.Params
				}
				if row.Name == default_connections {
					row.Open(Sqlserver(dsn), getConfig(row.Log))
				} else {
					Resolver = dbresolver.Register(dbresolver.Config{
						Replicas: []gorm.Dialector{Sqlserver(dsn)},
						// sources/replicas 负载均衡策略
						Policy: dbresolver.RandomPolicy{},
					}, row.Name)
				}
				break
			}
		}

	}

	if len(connections) > 1 {
		err := Master.Use(Resolver)
		if err != nil {
			panic(err)
		}
	}

}

// getConfig
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

// Open connection
func (d *Db) Open(Dialector gorm.Dialector, opts gorm.Option) {
	var db, err = gorm.Open(Dialector, opts)
	if err != nil {
		alog.Panic(err.Error())
	}
	Master = db
}

// Use
func (d *Db) Use(plugin gorm.Plugin) {
	if err := Master.Use(plugin); err != nil {
		alog.Panic(err.Error())
	}

}

// Mysql connection
func Mysql(dsn string) gorm.Dialector {
	return mysql.New(mysql.Config{
		DSN:                       dsn, // DSN data source name
		DefaultStringSize:         256, // string 类型字段的默认长度
		DefaultDatetimePrecision:  &datetimePrecision,
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	})
}

// Postgres connection
func Postgres(dsn string) gorm.Dialector {
	return postgres.New(postgres.Config{
		DriverName:           "",
		DSN:                  dsn,
		PreferSimpleProtocol: true,
		WithoutReturning:     false,
		Conn:                 nil,
	})
}

// Sqlserver connection
func Sqlserver(dsn string) gorm.Dialector {
	return sqlserver.Open(dsn)
}

// Distributed
func (d *Db) Distributed(config dbresolver.Config, datas ...interface{}) *dbresolver.DBResolver {
	return dbresolver.Register(config, datas)
}

// Close 关闭数据库
func Close() {
	if Master != nil {
		var db, err = Master.DB()
		if err != nil {
			alog.Error(err.Error())
		}

		if db != nil {
			if err2 := db.Close(); err2 != nil {
				return
			}
		}
	}
}
