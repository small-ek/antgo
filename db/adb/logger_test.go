package adb

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gormlogger "gorm.io/gorm/logger"
)

func TestFormatSQLForLogRemovesPostgresIdentifierQuotes(t *testing.T) {
	sql := "SELECT *\nFROM \"public\".\"sys_menu\"\r\nWHERE id in (1,2)\tAND \"sys_menu\".\"deleted_at\" IS NULL"

	got := formatSQLForLog(sql)
	want := "SELECT * FROM public.sys_menu WHERE id in (1,2) AND sys_menu.deleted_at IS NULL"
	if got != want {
		t.Fatalf("formatSQLForLog() = %q, want %q", got, want)
	}
}

func TestFormatSQLForLogKeepsDoubleQuotesInsideStringValues(t *testing.T) {
	sql := "SELECT * FROM \"public\".\"sys_role\" WHERE name='role \"admin\"' AND note='it''s ok'"

	got := formatSQLForLog(sql)
	want := "SELECT * FROM public.sys_role WHERE name='role \"admin\"' AND note='it''s ok'"
	if got != want {
		t.Fatalf("formatSQLForLog() = %q, want %q", got, want)
	}
}

func TestTraceLogsOriginalSQLError(t *testing.T) {
	var output bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&output),
		zap.ErrorLevel,
	)
	logger := Logger{
		ZapLogger: zap.New(core),
		LogLevel:  gormlogger.Error,
	}
	wantErr := errors.New("ERROR: deadlock detected (SQLSTATE 40P01)")

	logger.Trace(context.Background(), time.Now(), func() (string, int64) {
		return "UPDATE public.loan_face_liveness SET status=1 WHERE id=140", 0
	}, wantErr)

	var entry struct {
		Message string `json:"msg"`
		Error   string `json:"error"`
	}
	if err := json.Unmarshal(output.Bytes(), &entry); err != nil {
		t.Fatalf("unmarshal sql log: %v", err)
	}
	if entry.Message != "sql_error" {
		t.Fatalf("log message = %q, want sql_error", entry.Message)
	}
	if entry.Error != wantErr.Error() {
		t.Fatalf("error field = %q, want %q", entry.Error, wantErr.Error())
	}
}
