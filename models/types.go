package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	b, err := json.Marshal(s)
	return string(b), err
}

func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = StringSlice{}
		return nil
	}
	var b []byte
	switch v := value.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into StringSlice", value)
	}
	if len(b) == 0 {
		*s = StringSlice{}
		return nil
	}
	if err := json.Unmarshal(b, s); err != nil {
		// Legacy rows hold a bare scalar (e.g. "noun") rather than a JSON array.
		*s = StringSlice{string(b)}
	}
	return nil
}
