package nulls

import (
	"fmt"
)

type Bool struct {
	Value bool
	Valid bool
}

func (nt *Bool) Scan(value interface{}) error {
	switch value.(type) {
	case bool:
		*nt = Bool{value.(bool), true}
	default:
		*nt = Bool{false, false}
	}

	if !nt.Valid {
		return handleError(value)
	}

	return nil
}

// MarshalJSON for Bool
func (nt *Bool) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}

	val := fmt.Sprintf("%t", nt.Value)
	return []byte(val), nil
}

// UnmarshalJSON for Bool
func (nt *Bool) UnmarshalJSON(b []byte) error {
	s := string(b)

	if s != "true" && s != "false" {
		nt.Valid = false
		return nil
	}

	nt.Value = s == "true"
	nt.Valid = true

	return nil
}
