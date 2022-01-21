package sql

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/small-ek/antgo/utils/str"
	"time"
)

//Time 时间类型
type Time struct {
	sql.NullTime
}

// MarshalJSON 序列化为JSON
func (t Time) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("\"\""), nil
	}
	stamp := fmt.Sprintf("\"%s\"", t.Time.Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

// UnmarshalJSON 反序列化为JSON
func (t *Time) UnmarshalJSON(data []byte) error {
	times := str.ClearQuotes(string(data))
	if times != "" {
		var err error
		t.Time, err = time.Parse("2006-01-02 15:04:05", times)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

//Scan 读取插入
func (t *Time) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Time{sql.NullTime{Time: value}}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

//Value 插入
func (t Time) Value() (driver.Value, error) {
	strs := t.Time.Format("2006-01-02 15:04:05")
	if strs == "0001-01-01 00:00:00" {
		return nil, nil
	}
	return strs, nil
}
