package null

import (
	"database/sql/driver"
	"encoding/json"
)

// String represents a string that may be null.
// String implements the Scanner interface so
// it can be used as a scan destination.
// String implements the Marshaller interface
// so it can be used to read and write null
// values.
type String struct {
	String string
	Valid  bool // Valid is true if String is not ""
}

func (nt *String) Scan(value interface{}) error {
	nt.String, nt.Valid = value.(string)
	if nt.String == "" {
		nt.Valid = false
	}
	return nil
}

func (nt String) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.String, nil
}

func (nt String) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(nt.String)
}

func (nt *String) UnmarshalJSON(data []byte) error {
	if data == nil {
		nt.String = ""
	} else if err := json.Unmarshal(data, &nt.String); err != nil {
		return err
	}
	nt.Valid = nt.String != ""
	return nil
}

func NewString(t string) String {
	return String{
		String: t,
		Valid:  t != "",
	}
}
