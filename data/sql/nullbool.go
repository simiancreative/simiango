package sql

import (
	"fmt"
	"reflect"
)

type NullBool struct {
	Value bool
	Valid bool
}

func (nt *NullBool) Scan(value interface{}) error {
	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*nt = NullBool{nt.Value, false}
	} else {
		*nt = NullBool{value.(bool), true}
	}

	return nil
}

// MarshalJSON for NullBool
func (nt *NullBool) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}

	val := fmt.Sprintf("%t", nt.Value)
	return []byte(val), nil
}

// UnmarshalJSON for NullBool
func (nt *NullBool) UnmarshalJSON(b []byte) error {
	s := string(b)

	if s != "true" && s != "false" {
		nt.Valid = false
		return nil
	}

	nt.Value = s == "true"
	nt.Valid = true

	return nil
}
