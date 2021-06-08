package atime

import (
	"bytes"
	"os"
	"time"
)

// UTCS 将当前时间转换为 UTC 时区。
func (t *Times) UTCS() *Times {
	return WithTime(t.Time.UTC())
}

//Timestamp Get timestamp.时间转时间戳
func (t *Times) Timestamp() int64 {
	return t.Unix()
}

//Millisecond 毫秒
func (t *Times) Millisecond() int64 {
	return t.UnixNano() / 1e6
}

// Microsecond 微妙
func (t *Times) Microsecond() int64 {
	return t.UnixNano() / 1e3
}

// Nanosecond 纳秒
func (t *Times) Nanosecond() int64 {
	return t.UnixNano()
}

//Adds returns the time t+d.
func (t *Times) Adds(d time.Duration) *Times {
	timeStr := t.Add(d)
	return WithTime(timeStr)
}

// Months 返回 t 指定的年份中的月份。
func (t *Times) Months() int {
	return int(t.Month())
}

// Seconds 指定的分钟内的第二个偏移量
// 在 [0, 59] 范围内。
func (t *Times) Seconds() int {
	return t.Second()
}

// IsZeros reports whether t represents the zero time instant,
// January 1, year 1, 00:00:00 UTC.
func (t *Times) IsZeros() bool {
	if t == nil {
		return true
	}
	return t.IsZero()
}

// String returns current time object as string.
func (t *Times) String() string {
	if t == nil {
		return ""
	}
	if t.IsZero() {
		return ""
	}
	return t.Format("yyyy-MM-dd HH:mm:ss")
}

// AddDates adds year, month and day to the time.
func (t *Times) AddDates(years int, months int, days int) *Times {
	return WithTime(t.AddDate(years, months, days))
}

//WithDate ...
func WithDate(year, month, date, hour, minute, second int) *Times {
	t := time.Date(year, time.Month(month), date, hour, minute, second, 0, time.Local)
	return WithTime(t)
}

//SetTimeZone creates and returns a Time object with given timestamp
//Set time zone 设置时区
func SetTimeZone(zone string) error {
	location, err := time.LoadLocation(zone)
	if err != nil {
		return err
	}
	return os.Setenv("TZ", location.String())
}

// Truncates returns the result of rounding t down to a multiple of d (since the zero time).
// If d <= 0, Truncate returns t stripped of any monotonic clock reading but otherwise unchanged.
//
// Truncate operates on the time as an absolute duration since the
// zero time; it does not operate on the presentation form of the
// time. Thus, Truncate(Hour) may return a time with a non-zero
// minute, depending on the time's Location.
func (t *Times) Truncates(d time.Duration) *Times {
	return WithTime(t.Truncate(d))
}

// Equals reports whether t and u represent the same time instant.
// Two times can be equal even if they are in different locations.
// For example, 6:00 +0200 CEST and 4:00 UTC are Equal.
// See the documentation on the Time type for the pitfalls of using == with
// Time values; most code should use Equal instead.
func (t *Times) Equals(u *Times) bool {
	return t.Equal(u.Time)
}

// Befores reports whether the time instant t is before u.
func (t *Times) Befores(u *Times) bool {
	return t.Before(u.Time)
}

// Afters reports whether the time instant t is after u.
func (t *Times) Afters(u *Times) bool {
	return t.After(u.Time)
}

// Subs returns the duration t-u. If the result exceeds the maximum (or minimum)
// value that can be stored in a Duration, the maximum (or minimum) duration
// will be returned.
// To compute t-d for a duration d, use t.Add(-d).
func (t *Times) Subs(u *Times) time.Duration {
	return t.Sub(u.Time)
}

// StartOfMinute clones and returns a new time of which the seconds is set to 0.
func (t *Times) StartOfMinute() *Times {
	return WithTime(t.Truncate(time.Minute))
}

// StartOfHour clones and returns a new time of which the hour, minutes and seconds are set to 0.
func (t *Times) StartOfHour() *Times {
	y, m, d := t.Date()
	return WithTime(time.Date(y, m, d, t.Hour(), 0, 0, 0, t.Location()))
}

// StartOfDay clones and returns a new time which is the start of day, its time is set to 00:00:00.
func (t *Times) StartOfDay() *Times {
	y, m, d := t.Date()
	return WithTime(time.Date(y, m, d, 0, 0, 0, 0, t.Location()))
}

// StartOfWeek clones and returns a new time which is the first day of week and its time is set to
// 00:00:00.
func (t *Times) StartOfWeek() *Times {
	weekday := int(t.Weekday())
	return t.StartOfDay().AddDates(0, 0, -weekday)
}

// StartOfMonth clones and returns a new time which is the first day of the month and its is set to
// 00:00:00
func (t *Times) StartOfMonth() *Times {
	y, m, _ := t.Date()
	return WithTime(time.Date(y, m, 1, 0, 0, 0, 0, t.Location()))
}

// StartOfQuarter clones and returns a new time which is the first day of the quarter and its time is set
// to 00:00:00.
func (t *Times) StartOfQuarter() *Times {
	month := t.StartOfMonth()
	offset := (int(month.Month()) - 1) % 3
	return month.AddDates(0, -offset, 0)
}

// StartOfHalf clones and returns a new time which is the first day of the half year and its time is set
// to 00:00:00.
func (t *Times) StartOfHalf() *Times {
	month := t.StartOfMonth()
	offset := (int(month.Month()) - 1) % 6
	return month.AddDates(0, -offset, 0)
}

// StartOfYear clones and returns a new time which is the first day of the year and its time is set to
// 00:00:00.
func (t *Times) StartOfYear() *Times {
	y, _, _ := t.Date()
	return WithTime(time.Date(y, time.January, 1, 0, 0, 0, 0, t.Location()))
}

// EndOfMinute clones and returns a new time of which the seconds is set to 59.
func (t *Times) EndOfMinute() *Times {
	return t.StartOfMinute().Adds(time.Minute - time.Nanosecond)
}

// EndOfHour clones and returns a new time of which the minutes and seconds are both set to 59.
func (t *Times) EndOfHour() *Times {
	return t.StartOfHour().Adds(time.Hour - time.Nanosecond)
}

// EndOfDay clones and returns a new time which is the end of day the and its time is set to 23:59:59.
func (t *Times) EndOfDay() *Times {
	y, m, d := t.Date()
	return WithTime(time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location()))
}

// EndOfWeek clones and returns a new time which is the end of week and its time is set to 23:59:59.
func (t *Times) EndOfWeek() *Times {
	return t.StartOfWeek().AddDates(0, 0, 7).Adds(-time.Nanosecond)
}

// EndOfMonth clones and returns a new time which is the end of the month and its time is set to 23:59:59.
func (t *Times) EndOfMonth() *Times {
	return t.StartOfMonth().AddDates(0, 1, 0).Adds(-time.Nanosecond)
}

// EndOfQuarter clones and returns a new time which is end of the quarter and its time is set to 23:59:59.
func (t *Times) EndOfQuarter() *Times {
	return t.StartOfQuarter().AddDates(0, 3, 0).Adds(-time.Nanosecond)
}

// EndOfHalf clones and returns a new time which is the end of the half year and its time is set to 23:59:59.
func (t *Times) EndOfHalf() *Times {
	return t.StartOfHalf().AddDates(0, 6, 0).Adds(-time.Nanosecond)
}

// EndOfYear clones and returns a new time which is the end of the year and its time is set to 23:59:59.
func (t *Times) EndOfYear() *Times {
	return t.StartOfYear().AddDates(1, 0, 0).Adds(-time.Nanosecond)
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (t *Times) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (t *Times) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		t.Time = time.Time{}
		return nil
	}
	newTime := StrToTime(string(bytes.Trim(b, `"`)))
	t.Time = newTime.Time
	return nil
}
