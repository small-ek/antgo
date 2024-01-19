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
)

var Master map[string]*gorm.DB

// var Resolver *dbresolver.DBResolver
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
	if Master == nil {
		Master = make(map[string]*gorm.DB)
	}
	connections := config.GetMaps("connections")
	default_connections := conv.String(config.Get("system.default_connections"))
	if default_connections != "" {

		for i := 0; i < len(connections); i++ {
			value := connections[i]
			row := Db{}
			conv.Struct(&row, value)
			dsn := row.Dsn

			switch row.Type {
			case "mysql":
				if row.Dsn == "" && row.Name != "" {
					dsn = row.Username + ":" + row.Password + "@tcp(" + row.Hostname + ":" + row.Port + ")/" + row.Database + "?" + row.Params
					Master[row.Name], _ = row.Open(Mysql(dsn), getConfig(row.Log))
				}
				break
			case "pgsql":
				if row.Dsn == "" && row.Name != "" {
					dsn = "host=" + row.Hostname + " port=" + row.Port + " user=" + row.Username + " dbname=" + row.Database + " " + row.Params + " password=" + row.Password + row.Params
					Master[row.Name], _ = row.Open(Postgres(dsn), getConfig(row.Log))
				}
				break
			case "sqlsrv":
				if row.Dsn == "" && row.Name != "" {
					dsn = "sqlserver://" + row.Username + ":" + row.Password + "@" + row.Hostname + ":" + row.Port + "?database=" + row.Database + row.Params
					Master[row.Name], _ = row.Open(Sqlserver(dsn), getConfig(row.Log))
				}
				break
			}
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
func (d *Db) Open(Dialector gorm.Dialector, opts gorm.Option) (db *gorm.DB, err error) {
	db, err = gorm.Open(Dialector, opts)
	if err != nil {
		alog.Panic("Open", zap.Error(err))
	}
	return
}

// Use
func (d *Db) Use(name string, plugin gorm.Plugin) {
	connections := config.GetMaps("connections")

	for i := 0; i < len(connections); i++ {
		value := connections[i]
		row := Db{}
		conv.Struct(&row, value)
		if name == row.Name {
			if err := Master[row.Name].Use(plugin); err != nil {
				alog.Error("Use", zap.Error(err))
			}
		}

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
	connections := config.GetMaps("connections")

	for i := 0; i < len(connections); i++ {
		value := connections[i]
		row := Db{}
		conv.Struct(&row, value)
		if Master[row.Name] != nil {
			var db, err = Master[row.Name].DB()
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

}
