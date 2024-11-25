package adb

import (
	"database/sql"
	"fmt"
	"github.com/small-ek/antgo/os/alog"
	"github.com/small-ek/antgo/utils/conv"
	"go.uber.org/zap"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"sync"
	"time"
)

var (
	Master map[string]*gorm.DB
	once   sync.Once
)

type DatabaseConfig struct {
	Name            string              `json:"name"`
	Type            string              `json:"type"`
	Hostname        string              `json:"hostname"`
	Port            string              `json:"port"`
	Username        string              `json:"username"`
	Password        string              `json:"password"`
	DatabaseName    string              `json:"database"`
	Params          string              `json:"params"`
	LogEnabled      bool                `json:"log"`
	DSN             string              `json:"dsn"`
	MaxIdleConns    int                 `json:"max_idle_conns"`
	MaxOpenConns    int                 `json:"max_open_conns"`
	ConnMaxLifetime int                 `json:"conn_max_lifetime"`
	ConnMaxIdleTime int                 `json:"conn_max_idleTime"`
	LogLevel        gormlogger.LogLevel `json:"level"`
}

// InitDb initializes the database connections based on the configuration 初始化基于配置的数据库连接
func InitDb(connections []map[string]any) {
	once.Do(func() {
		Master = make(map[string]*gorm.DB)

		for i := 0; i < len(connections); i++ {
			value := connections[i]
			config := DatabaseConfig{}
			conv.ToStruct(value, &config)

			if value["name"] != "" {
				db, err := CreateConnection(config)
				if err != nil {
					alog.Write.Panic("Failed to initialize database connection "+config.Name+" :", zap.Error(err))
				}
				Master[config.Name] = db
			}

		}
	})

}

// GetDatabase retrieves a database connection by name 通过名称检索数据库连接
func GetDatabase(name string) *gorm.DB {
	db, exists := Master[name]
	if !exists {
		alog.Write.Error("Database connection not found:", zap.String("name", name))
	}
	return db
}

// generateDSN generates a DSN string based on the database config 生成基于数据库配置的DSN字符串
func generateDSN(config DatabaseConfig) string {
	switch config.Type {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
			config.Username, config.Password, config.Hostname, config.Port, config.DatabaseName, config.Params)
	case "pgsql":
		return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s %s",
			config.Hostname, config.Port, config.Username, config.DatabaseName, config.Password, config.Params)
	case "sqlsrv":
		return fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&%s",
			config.Username, config.Password, config.Hostname, config.Port, config.DatabaseName, config.Params)
	case "clickhouse":
		return fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s?%s",
			config.Username, config.Password, config.Hostname, config.Port, config.DatabaseName, config.Params)
	default:
		return ""
	}
}

// configureConnectionPool sets up the connection pool parameters for the database connection 设置数据库连接的连接池参数
func configureConnectionPool(sqlDB *sql.DB, config DatabaseConfig) {
	if config.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)
	}
	if config.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(time.Duration(config.ConnMaxIdleTime) * time.Hour)
	}
}

// CreateConnection creates a single database connection based on the configuration 创建基于配置的单个数据库连接
func CreateConnection(config DatabaseConfig) (db *gorm.DB, err error) {
	if config.DSN == "" {
		config.DSN = generateDSN(config)
	}

	switch config.Type {
	case "mysql":
		db, err = config.Open(Mysql(config.DSN), getConfig(config.LogEnabled, config.LogLevel))
	case "pgsql":
		db, err = config.Open(Postgres(config.DSN), getConfig(config.LogEnabled, config.LogLevel))
	case "sqlsrv":
		db, err = config.Open(Sqlserver(config.DSN), getConfig(config.LogEnabled, config.LogLevel))
	case "clickhouse":
		db, err = config.Open(clickhouse.Open(config.DSN), getConfig(config.LogEnabled, config.LogLevel))
	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.Type)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	configureConnectionPool(sqlDB, config)
	return db, nil
}

// getConfig
func getConfig(isLog bool, level gormlogger.LogLevel) *gorm.Config {
	if isLog {
		zapLog := New(zap.L())
		zapLog.SetAsDefault()
		return &gorm.Config{
			Logger:                                   zapLog.LogMode(level),
			DisableForeignKeyConstraintWhenMigrating: true,
			SkipDefaultTransaction:                   true,
			PrepareStmt:                              true,
			AllowGlobalUpdate:                        false,
		}
	} else {
		return &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			SkipDefaultTransaction:                   true,
			PrepareStmt:                              true,
			AllowGlobalUpdate:                        false,
		}
	}
}

// Open connection
func (d *DatabaseConfig) Open(Dialector gorm.Dialector, opts gorm.Option) (db *gorm.DB, err error) {
	db, err = gorm.Open(Dialector, opts)
	if err != nil {
		alog.Write.Panic("gorm open error:", zap.Error(err))
	}
	return
}

// Use
func (d *DatabaseConfig) Use(name string, plugin gorm.Plugin) {
	if err := Master[name].Use(plugin); err != nil {
		alog.Write.Error("Use", zap.Error(err))
	}
}

// Mysql connection Mysql连接
func Mysql(dsn string) gorm.Dialector {
	defaultPrecision := 3 // 设置默认的时间精度为 3（毫秒）
	return mysql.New(mysql.Config{
		DSN:                       dsn, // DSN data source name
		DefaultStringSize:         256, // string 类型字段的默认长度
		DefaultDatetimePrecision:  &defaultPrecision,
		DisableDatetimePrecision:  false, // 支持 datetime 精度
		DontSupportRenameIndex:    false, // 支持直接重命名索引（适用于 MySQL 8.0+）
		DontSupportRenameColumn:   false, // 支持直接重命名列（适用于 MySQL 8.0+）
		SkipInitializeWithVersion: false, // 自动初始化版本特性
	})
}

// Postgres connection Postgres连接
func Postgres(dsn string) gorm.Dialector {
	return postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
		WithoutReturning:     false,
	})
}

// Sqlserver connection Sqlserver连接
func Sqlserver(dsn string) gorm.Dialector {
	return sqlserver.Open(dsn)
}

// Distributed 分布式
func (d *DatabaseConfig) Distributed(config dbresolver.Config, datas ...interface{}) *dbresolver.DBResolver {
	return dbresolver.Register(config, datas)
}

// Close 关闭数据库
func Close() {
	for name, db := range Master {
		if name != "default" {
			sqlDB, err := db.DB()
			if err != nil {
				alog.Write.Error("Error retrieving DB instance for "+name+":", zap.Error(err))
				continue
			}
			if err := sqlDB.Close(); err != nil {
				alog.Write.Error("Error closing DB connection for "+name+":", zap.Error(err))
			} else {
				alog.Write.Warn("Database connection '" + name + "' closed successfully")
			}
		}

	}

}
