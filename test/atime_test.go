package test

import (
	"github.com/small-ek/antgo/os/atime"
	"testing"
	"time"
)

// TestUTC 测试 UTC 方法
func TestUTC(t *testing.T) {
	localTime := time.Date(2023, 10, 5, 12, 0, 0, 0, time.Local)
	times := &atime.Times{Time: localTime}

	utcTime := times.UTC()
	expected := localTime.UTC()

	if !utcTime.Time.Equal(expected) {
		t.Errorf("UTC() failed, expected %v, got %v", expected, utcTime.Time)
	}
}

// TestTimestamp 测试 Timestamp 方法
func TestTimestamp(t *testing.T) {
	now := time.Now()
	times := &atime.Times{Time: now}

	timestamp := times.Timestamp()
	expected := now.Unix()

	if timestamp != expected {
		t.Errorf("Timestamp() failed, expected %d, got %d", expected, timestamp)
	}
}

// TestAdd 测试 Add 方法
func TestAdd(t *testing.T) {
	now := time.Now()
	times := &atime.Times{Time: now}

	duration := time.Hour * 2
	newTime := times.Add(duration)
	expected := now.Add(duration)

	if !newTime.Time.Equal(expected) {
		t.Errorf("Add() failed, expected %v, got %v", expected, newTime.Time)
	}
}

// TestStartOfDay 测试 StartOfDay 方法
func TestStartOfDay(t *testing.T) {
	now := time.Date(2023, 10, 5, 15, 30, 45, 0, time.Local)
	times := &atime.Times{Time: now}

	startOfDay := times.StartOfDay()
	expected := time.Date(2023, 10, 5, 0, 0, 0, 0, time.Local)

	if !startOfDay.Time.Equal(expected) {
		t.Errorf("StartOfDay() failed, expected %v, got %v", expected, startOfDay.Time)
	}
}

// TestFormat 测试 Format 方法
func TestFormat(t *testing.T) {
	now := time.Date(2023, 10, 5, 15, 30, 45, 0, time.Local)
	times := &atime.Times{Time: now}

	formatted := times.Format("yyyy-MM-dd HH:mm:ss")
	expected := "2023-10-05 15:30:45"

	if formatted != expected {
		t.Errorf("Format() failed, expected %s, got %s", expected, formatted)
	}
}

// TestStrToTime 测试 StrToTime 方法
func TestStrToTime(t *testing.T) {
	str := "2023-10-05 15:30:45"
	times := atime.StrToTime(str)
	expected := time.Date(2023, 10, 5, 15, 30, 45, 0, time.Local)

	if !times.Time.Equal(expected) {
		t.Errorf("StrToTime() failed, expected %v, got %v", expected, times.Time)
	}
}

// TestNewFromTimeStamp 测试 NewFromTimeStamp 方法
func TestNewFromTimeStamp(t *testing.T) {
	timestamp := int64(1696527045) // 对应 2023-10-05 15:30:45 UTC
	times := atime.NewFromTimeStamp(timestamp)
	expected := time.Unix(timestamp, 0)

	if !times.Time.Equal(expected) {
		t.Errorf("NewFromTimeStamp() failed, expected %v, got %v", expected, times.Time)
	}
}

// TestEndOfMonth 测试 EndOfMonth 方法
func TestEndOfMonth(t *testing.T) {
	now := time.Date(2023, 10, 5, 15, 30, 45, 0, time.Local)
	times := &atime.Times{Time: now}

	endOfMonth := times.EndOfMonth()
	expected := time.Date(2023, 10, 31, 23, 59, 59, int(time.Second-time.Nanosecond), time.Local)

	if !endOfMonth.Time.Equal(expected) {
		t.Errorf("EndOfMonth() failed, expected %v, got %v", expected, endOfMonth.Time)
	}
}
