package nulls

import (
	"encoding/json"
)

type BitToInt64 struct {
	Value int64
	Valid bool
}

func (nt *BitToInt64) Scan(value interface{}) error {
	var bit int64
	bit = 0

	switch value.(type) {
	case bool:
		if value.(bool) {
			bit = 1
		}
		*nt = BitToInt64{bit, true}
	case nil:
		*nt = BitToInt64{int64(0), false}
	default:
		return handleError(value)
	}

	return nil
}

// MarshalJSON for Bool
func (nt *BitToInt64) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(nt.Value)
}

// UnmarshalJSON for Bool
func (nt *BitToInt64) UnmarshalJSON(b []byte) error {
	s := string(b)

	if s != "true" && s != "false" {
		nt.Valid = false
		return nil
	}

	nt.Value = 0

	if s == "true" {
		nt.Value = 1
	}

	nt.Valid = true

	return nil
}
