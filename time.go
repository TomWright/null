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
	time.Time
}

// NewTime returns a Time containing the given time.Time.
func NewTime(t time.Time) Time {
	return Time{
		Time: t,
	}
}

// Valid returns true if the time.Time contained within Time is
// a non-zero time object.
func (nt *Time) Valid() bool {
	return !nt.IsZero()
}

// Scan implements the Scanner interface.
func (nt *Time) Scan(value interface{}) error {
	var valid bool
	nt.Time, valid = value.(time.Time)
	if !valid {
		nt.Time = time.Time{}
	}
	return nil
}

// Value implements the driver Valuer interface.
func (nt Time) Value() (driver.Value, error) {
	if !nt.Valid() {
		return nil, nil
	}
	return nt.Time, nil
}

// MarshalJSON returns either a marshal'd time.Time if it was
// valid, or a marshal'd NULL value if it was not valid.
func (nt Time) MarshalJSON() ([]byte, error) {
	if !nt.Valid() {
		return json.Marshal(nil)
	}
	return json.Marshal(nt.Time)
}

// UnmarshalJSON attempts to unmarshal the given bytes
// into a Time object.
func (nt *Time) UnmarshalJSON(data []byte) error {
	var t time.Time
	if data != nil {
		if err := json.Unmarshal(data, &t); err != nil {
			return err
		}
	}
	nt.Time = t
	return nil
}
