package adb

import (
	"testing"

	gormlogger "gorm.io/gorm/logger"
)

func TestGetConfigPrepareStmt(t *testing.T) {
	if cfg := getConfig(false, gormlogger.Silent, false); cfg.PrepareStmt {
		t.Fatal("PostgreSQL configuration must disable prepared statements")
	}
	if cfg := getConfig(false, gormlogger.Silent, true); !cfg.PrepareStmt {
		t.Fatal("non-PostgreSQL configuration must preserve prepared statements")
	}
}
