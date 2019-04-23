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
func (nt Time) Valid() bool {
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

// Override some embedded functions taken from time.Time to make
// things easier to use.

// Add returns the time t+d.
func (nt Time) Add(d time.Duration) Time {
	return NewTime(nt.Time.Add(d))
}

// UTC returns nt with the location set to UTC.
func (nt Time) UTC() Time {
	return NewTime(nt.Time.UTC())
}

// AddDate returns the time corresponding to adding the
// given number of years, months, and days to nt.
// For example, AddDate(-1, 2, 3) applied to January 1, 2011
// returns March 4, 2010.
//
// AddDate normalizes its result in the same way that Date does,
// so, for example, adding one month to October 31 yields
// December 1, the normalized form for November 31.
func (nt Time) AddDate(years int, months int, days int) Time {
	return NewTime(nt.Time.AddDate(years, months, days))
}

// In returns a copy of nt representing the same time instant, but
// with the copy's location information set to loc for display
// purposes.
//
// In panics if loc is nil.
func (nt Time) In(loc *time.Location) Time {
	return NewTime(nt.Time.In(loc))
}

// Local returns nt with the location set to local time.
func (nt Time) Local(d time.Duration) Time {
	return NewTime(nt.Time.Local())
}

// Round returns the result of rounding nt to the nearest multiple of d (since the zero time).
// The rounding behavior for halfway values is to round up.
// If d <= 0, Round returns nt stripped of any monotonic clock reading but otherwise unchanged.
//
// Round operates on the time as an absolute duration since the
// zero time; it does not operate on the presentation form of the
// time. Thus, Round(Hour) may return a time with a non-zero
// minute, depending on the time's Location.
func (nt Time) Round(d time.Duration) Time {
	return NewTime(nt.Time.Round(d))
}

// Truncate returns the result of rounding nt down to a multiple of d (since the zero time).
// If d <= 0, Truncate returns nt stripped of any monotonic clock reading but otherwise unchanged.
//
// Truncate operates on the time as an absolute duration since the
// zero time; it does not operate on the presentation form of the
// time. Thus, Truncate(Hour) may return a time with a non-zero
// minute, depending on the time's Location.
func (nt Time) Truncate(d time.Duration) Time {
	return NewTime(nt.Time.Truncate(d))
}
