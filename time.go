package null

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Time implements the Scanner interface so
// it can be used as a scan destination.
// Time implements the Marshaller interface
// so it can be used to read and write null
// values.
type Time struct {
	Time  time.Time
	Valid bool
}

// Scan implements the Scanner interface.
func (nt *Time) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	if nt.Valid && nt.Time.IsZero() {
		nt.Valid = false
	}
	return nil
}

// Value implements the driver Valuer interface.
func (nt Time) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// TimeOrZero returns either a valid time.Time, or a zero-value
// time.Time object.
func (nt Time) TimeOrZero() time.Time {
	if !nt.Valid {
		return time.Time{}
	}
	return nt.Time
}

// MarshalJSON returns either a marshal'd time.Time if it was
// valid, or a marshal'd NULL value if it was not valid.
func (nt Time) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(nt.Time)
}

func (nt *Time) UnmarshalJSON(data []byte) error {
	var t time.Time
	if data != nil {
		if err := json.Unmarshal(data, &t); err != nil {
			return err
		}
	}
	nt.Time = t
	nt.Valid = !t.IsZero()
	return nil
}

func NewTime(t time.Time) Time {
	return Time{
		Time:  t,
		Valid: !t.IsZero(),
	}
}
