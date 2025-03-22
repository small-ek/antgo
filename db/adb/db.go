package adb

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

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
)

// Global Master map 存储所有数据库连接对象（全局唯一）
var (
	Master = make(map[string]*gorm.DB)
	once   sync.Once
)

// DatabaseConfig 数据库配置结构体
// DatabaseConfig represents the configuration for a database connection.
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
	ConnMaxIdleTime int                 `json:"conn_max_idle_time"` // 注意：修改为统一格式
	LogLevel        gormlogger.LogLevel `json:"level"`
}

// InitDb 初始化数据库连接
// InitDb initializes database connections based on the provided configuration.
// connections is a slice of maps representing the configuration for each database.
func InitDb(connections []map[string]any) {
	once.Do(func() {
		for _, value := range connections {
			var config DatabaseConfig
			if err := conv.ToStruct(value, &config); err != nil {
				alog.Write.Panic("Failed to convert database config to struct:", zap.Error(err))
			}

			// 仅当 name 不为空时初始化连接 / Initialize connection only if name is provided
			if config.Name != "" {
				db, err := CreateConnection(config)
				if err != nil {
					alog.Write.Panic("Failed to initialize database connection "+config.Name+" : ", zap.Error(err))
				}
				Master[config.Name] = db
			}
		}
	})
}

// GetDatabase 根据名称获取数据库连接对象
// GetDatabase retrieves a database connection by its name.
func GetDatabase(name string) *gorm.DB {
	db, exists := Master[name]
	if !exists {
		alog.Write.Error("Database connection not found :", zap.String("name", name))
	}
	return db
}

// generateDSN 根据数据库配置生成 DSN 字符串
// generateDSN generates a DSN string based on the provided database configuration.
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

// configureConnectionPool 设置数据库连接池参数
// configureConnectionPool configures the connection pool settings for a database connection.
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

// CreateConnection 根据配置创建单个数据库连接
// CreateConnection creates a single database connection based on the provided configuration.
func CreateConnection(config DatabaseConfig) (*gorm.DB, error) {
	// 若未提供 DSN，则自动生成 / Automatically generate DSN if not provided.
	if config.DSN == "" {
		config.DSN = generateDSN(config)
	}

	var (
		db  *gorm.DB
		err error
	)
	// 根据数据库类型选择对应驱动 / Select appropriate driver based on database type.
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

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	configureConnectionPool(sqlDB, config)
	return db, nil
}

// getConfig 返回 GORM 的配置对象
// getConfig returns a GORM configuration object based on logging settings.
func getConfig(isLog bool, level gormlogger.LogLevel) *gorm.Config {
	cfg := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		PrepareStmt:                              true,
		AllowGlobalUpdate:                        false,
	}

	if isLog {
		// 初始化自定义 Zap 日志记录器，并设置为默认日志记录器 / Initialize a custom Zap logger and set it as the default logger.
		zapLog := New(alog.Write)
		zapLog.SetAsDefault()
		cfg.Logger = zapLog.LogMode(level)
	}
	return cfg
}

// Open 使用指定的 Dialector 和选项打开数据库连接
// Open opens a database connection using the provided Dialector and GORM options.
func (d *DatabaseConfig) Open(dialector gorm.Dialector, opts *gorm.Config) (*gorm.DB, error) {
	db, err := gorm.Open(dialector, opts)
	if err != nil {
		// 记录错误，但不直接 panic，交由调用者处理 / Log error and return it for the caller to handle.
		alog.Write.Panic("gorm open error :", zap.Error(err))
	}
	return db, nil
}

// Use 为指定数据库连接添加插件
// Use applies a GORM plugin to the database connection identified by name.
func (d *DatabaseConfig) Use(name string, plugin gorm.Plugin) {
	if db, exists := Master[name]; exists {
		if err := db.Use(plugin); err != nil {
			alog.Write.Error("Failed to apply plugin:", zap.Error(err))
		}
	} else {
		alog.Write.Error("Database connection not found when applying plugin:", zap.String("name", name))
	}
}

// Mysql 返回 MySQL 的 Dialector 配置
// Mysql returns the GORM Dialector for MySQL.
func Mysql(dsn string) gorm.Dialector {
	defaultPrecision := 3 // 默认时间精度设为3（毫秒） / Set default time precision to 3 (milliseconds)
	return mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256, // 默认字符串长度 / Default length for string fields
		DefaultDatetimePrecision:  &defaultPrecision,
		DisableDatetimePrecision:  false, // 支持 datetime 精度 / Enable datetime precision
		DontSupportRenameIndex:    false, // 支持直接重命名索引（适用于 MySQL 8.0+） / Support renaming indexes (MySQL 8.0+)
		DontSupportRenameColumn:   false, // 支持直接重命名列（适用于 MySQL 8.0+） / Support renaming columns (MySQL 8.0+)
		SkipInitializeWithVersion: false, // 自动初始化版本特性 / Automatically initialize version-specific features
	})
}

// Postgres 返回 PostgreSQL 的 Dialector 配置
// Postgres returns the GORM Dialector for PostgreSQL.
func Postgres(dsn string) gorm.Dialector {
	return postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
		WithoutReturning:     false,
	})
}

// Sqlserver 返回 SQL Server 的 Dialector 配置
// Sqlserver returns the GORM Dialector for SQL Server.
func Sqlserver(dsn string) gorm.Dialector {
	return sqlserver.Open(dsn)
}

// Distributed 设置分布式数据库，返回 dbresolver 插件实例
// Distributed configures a distributed database using GORM's dbresolver plugin.
// 'config' specifies the resolver configuration, and 'datas' represents the data sources.
func (d *DatabaseConfig) Distributed(config dbresolver.Config, datas ...interface{}) *dbresolver.DBResolver {
	return dbresolver.Register(config, datas)
}

// Close 关闭所有数据库连接（除了名称为 "default" 的连接）
// Close closes all database connections except for the one named "default".
func Close() {
	for name, db := range Master {
		if name == "default" {
			continue
		}
		sqlDB, err := db.DB()
		if err != nil {
			alog.Write.Error("Error retrieving DB instance for "+name+"：", zap.Error(err))
			continue
		}
		if err = sqlDB.Close(); err != nil {
			alog.Write.Error("Error closing DB connection for "+name+"：", zap.Error(err))
		} else {
			alog.Write.Info("Database connection '" + name + "' closed successfully")
		}
	}
}
