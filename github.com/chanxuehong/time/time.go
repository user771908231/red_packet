package time

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	sqlx "github.com/chanxuehong/database/sql"
)

type Time struct {
	time.Time
}

func Date(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) Time {
	return Time{Time: time.Date(year, month, day, hour, min, sec, nsec, loc)}
}

func Now() Time { return Time{Time: time.Now()} }

func Parse(layout, value string) (t Time, err error) {
	tt, err := time.Parse(layout, value)
	if err != nil {
		return
	}
	t = Time{Time: tt}
	return
}

func ParseInLocation(layout, value string, loc *time.Location) (t Time, err error) {
	tt, err := time.ParseInLocation(layout, value, loc)
	if err != nil {
		return
	}
	t = Time{Time: tt}
	return
}

func Unix(sec int64, nsec int64) Time { return Time{Time: time.Unix(sec, nsec)} }

var _ driver.Valuer = Time{}
var _ sql.Scanner = (*Time)(nil)

// Value implements the driver.Valuer interface.
// It converts Time to unixtime, zero Time converts to 0.
func (t Time) Value() (value driver.Value, err error) {
	if t.IsZero() {
		value = 0
	} else {
		value = t.Unix()
	}
	return
}

var zeroTime time.Time

// Scan implements the sql.Scanner interface.
// It scans unixtime to Time, 0 means zero Time.
func (t *Time) Scan(value interface{}) (err error) {
	if value == nil {
		return errors.New("Can't convert nil value to Time")
	}
	var unixtime int64
	if err = sqlx.ConvertAssign(&unixtime, value); err != nil {
		return
	}
	if unixtime == 0 {
		t.Time = zeroTime
	} else {
		t.Time = time.Unix(unixtime, 0)
	}
	return
}

var _ json.Marshaler = Time{}
var _ json.Unmarshaler = (*Time)(nil)

// MarshalJSON implements the json.Marshaler interface.
// The time is a quoted string in "2006-01-02 15:04:05" format.
func (t Time) MarshalJSON() ([]byte, error) {
	if y := t.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}
	return []byte(t.In(BeijingLocation).Format(`"2006-01-02 15:04:05"`)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in "2006-01-02 15:04:05" format.
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	t.Time, err = time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), BeijingLocation)
	return
}
