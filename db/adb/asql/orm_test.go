package asql

import (
	"strings"
	"testing"

	"github.com/small-ek/antgo/utils/page"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newDryRunDB(t *testing.T, dialector gorm.Dialector) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(dialector, &gorm.Config{
		DryRun:               true,
		DisableAutomaticPing: true,
	})
	if err != nil {
		t.Fatalf("open dry-run db: %v", err)
	}
	return db
}

func newPostgresDryRunDB(t *testing.T) *gorm.DB {
	return newDryRunDB(t, postgres.New(postgres.Config{
		DriverName:           "pgx",
		DSN:                  "host=localhost user=gorm dbname=gorm password=gorm sslmode=disable",
		PreferSimpleProtocol: true,
	}))
}

func newMysqlDryRunDB(t *testing.T) *gorm.DB {
	return newDryRunDB(t, mysql.New(mysql.Config{
		DSN:                       "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local",
		SkipInitializeWithVersion: true,
	}))
}

func scopeSQL(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (string, []interface{}) {
	var rows []map[string]interface{}
	tx := db.Table("public.sys_api").Scopes(scopes...).Find(&rows)
	return tx.Statement.SQL.String(), tx.Statement.Vars
}

func TestWhereQuotesPostgresColumns(t *testing.T) {
	sql, vars := scopeSQL(newPostgresDryRunDB(t), Where("path", "ILIKE", "channel-new-user-rep"))

	if strings.Contains(sql, "`") {
		t.Fatalf("postgres sql contains mysql backticks: %s", sql)
	}
	if !strings.Contains(sql, `"path" ILIKE $1`) {
		t.Fatalf("postgres sql should quote column with double quotes: %s", sql)
	}
	if len(vars) != 1 || vars[0] != "channel-new-user-rep%" {
		t.Fatalf("unexpected vars: %#v", vars)
	}
}

func TestWhereKeepsMysqlQuoting(t *testing.T) {
	sql, vars := scopeSQL(newMysqlDryRunDB(t), Like("path", "channel-new-user-rep"))

	if !strings.Contains(sql, "`path` LIKE ?") {
		t.Fatalf("mysql sql should quote column with backticks: %s", sql)
	}
	if len(vars) != 1 || vars[0] != "channel-new-user-rep%" {
		t.Fatalf("unexpected vars: %#v", vars)
	}
}

func TestWhereSkipsZeroValues(t *testing.T) {
	sql, vars := scopeSQL(newPostgresDryRunDB(t),
		Where("id", "=", 0),
		Where("path", "LIKE", ""),
	)

	if strings.Contains(sql, " WHERE ") {
		t.Fatalf("zero values should not add where conditions: %s", sql)
	}
	if len(vars) != 0 {
		t.Fatalf("zero values should not bind vars: %#v", vars)
	}
}

func TestFiltersUsePostgresRegexOperator(t *testing.T) {
	sql, vars := scopeSQL(newPostgresDryRunDB(t), Filters([]page.Filter{
		{Field: "path", Operator: "RLIKE", Value: "^/api"},
	}))

	if strings.Contains(sql, "RLIKE") {
		t.Fatalf("postgres sql should not use mysql RLIKE: %s", sql)
	}
	if !strings.Contains(sql, `"path" ~ $1`) {
		t.Fatalf("postgres regex should use POSIX match operator: %s", sql)
	}
	if len(vars) != 1 || vars[0] != "^/api" {
		t.Fatalf("unexpected vars: %#v", vars)
	}
}
