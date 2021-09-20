package sql

import (
	"encoding/json"
	"reflect"
)

// NullString is an alias for sql.NullString data type
type NullString struct {
	Value string
	Valid bool
}

// Scan implements the Scanner interface for NullString
func (ns *NullString) Scan(value interface{}) error {
	if reflect.TypeOf(value) == nil {
		*ns = NullString{"", false}
		return nil
	}

	switch value.(type) {
	case string:
		*ns = NullString{value.(string), true}
	case []uint8:
		*ns = NullString{string(value.([]uint8)), true}
	default:
		*ns = NullString{"", false}
	}

	return nil
}

// MarshalJSON for NullString
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.Value)
}

// UnmarshalJSON for NullString
func (ns *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.Value)
	ns.Valid = (err == nil)
	return err
}
