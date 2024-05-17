package adb

import (
	"fmt"
	"github.com/small-ek/antgo/os/alog"
	"github.com/small-ek/antgo/utils/conv"
	"go.uber.org/zap"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"time"
)

var Master map[string]*gorm.DB

// var Resolver *dbresolver.DBResolver
var datetimePrecision = 2

type Db struct {
	Name            string `json:"name"`
	Type            string `json:"type"`
	Hostname        string `json:"hostname"`
	Port            string `json:"port"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	Database        string `json:"database"`
	Params          string `json:"params"`
	Log             bool   `json:"log"`
	Dsn             string `json:"dsn"`
	MaxIdleConns    int    `json:"max_idle_conns"`
	MaxOpenConns    int    `json:"max_open_conns"`
	ConnMaxLifetime int    `json:"conn_max_lifetime"`
}

// InitDb
func InitDb(connections []map[string]any) {
	if Master == nil {
		Master = make(map[string]*gorm.DB)
	}

	for i := 0; i < len(connections); i++ {
		value := connections[i]
		row := Db{}
		conv.Struct(&row, value)
		dsn := row.Dsn
		var err error
		if row.Name != "" {
			switch row.Type {
			case "mysql":
				if dsn == "" {
					dsn = row.Username + ":" + row.Password + "@tcp(" + row.Hostname + ":" + row.Port + ")/" + row.Database + "?" + row.Params
				}

				Master[row.Name], err = row.Open(Mysql(dsn), getConfig(row.Log))

				break
			case "pgsql":
				if dsn == "" {
					dsn = "host=" + row.Hostname + " port=" + row.Port + " user=" + row.Username + " dbname=" + row.Database + " " + row.Params + " password=" + row.Password + row.Params
				}

				Master[row.Name], err = row.Open(Postgres(dsn), getConfig(row.Log))

				break
			case "sqlsrv":
				if dsn == "" {
					dsn = "sqlserver://" + row.Username + ":" + row.Password + "@" + row.Hostname + ":" + row.Port + "?database=" + row.Database + row.Params
				}

				Master[row.Name], err = row.Open(Sqlserver(dsn), getConfig(row.Log))
				break
			case "clickhouse":
				if dsn == "" {
					dsn = "clickhouse://" + row.Username + ":" + row.Password + "@" + row.Hostname + ":" + row.Port + "/" + row.Database + row.Params
				}

				Master[row.Name], err = row.Open(clickhouse.Open(dsn), getConfig(row.Log))
				break
			}

			if err != nil {
				alog.Write.Panic("gorm open error:", zap.Error(err))
			}

			sqlDB, err := Master[row.Name].DB()
			if err != nil {
				alog.Write.Panic("gorm db error:", zap.Error(err))
			}
			//SetMaxIdleConns设置空闲连接池中的最大连接数。
			if row.MaxIdleConns > 0 {
				sqlDB.SetMaxIdleConns(row.MaxIdleConns)
			}

			// SetMaxOpenConns设置数据库的最大打开连接数。
			if row.MaxOpenConns > 0 {
				sqlDB.SetMaxOpenConns(row.MaxOpenConns)
			}

			// SetConnMaxLifetime设置连接可能被重用的最大时间。
			if row.ConnMaxLifetime > 0 {
				sqlDB.SetConnMaxLifetime(time.Duration(row.ConnMaxLifetime) * time.Second)
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
			SkipDefaultTransaction:                   true,
			PrepareStmt:                              true,
		}
	} else {
		return &gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
		}
	}
}

// Open connection
func (d *Db) Open(Dialector gorm.Dialector, opts gorm.Option) (db *gorm.DB, err error) {
	db, err = gorm.Open(Dialector, opts)
	if err != nil {
		alog.Write.Panic("gorm open error:", zap.Error(err))
	}
	return
}

// Use
func (d *Db) Use(name string, plugin gorm.Plugin) {
	if err := Master[name].Use(plugin); err != nil {
		alog.Write.Error("Use", zap.Error(err))
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
func Close(connections []map[string]any) {
	for i := 0; i < len(connections); i++ {
		value := connections[i]
		row := Db{}
		conv.Struct(&row, value)
		if Master[row.Name] != nil {
			var db, err = Master[row.Name].DB()
			if err != nil {
				alog.Write.Error("Close database", zap.Error(fmt.Errorf("failed to close database connection for %s: %s", row.Name, err.Error())))
			}

			if db != nil {
				if err2 := db.Close(); err2 != nil {
					alog.Write.Error("Close database", zap.Error(fmt.Errorf("failed to close database connection for %s: %s", row.Name, err2.Error())))
					return
				}
			}
		}
	}

}
