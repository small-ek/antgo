package adb

import "testing"

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
