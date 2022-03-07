package nulls

import (
	"encoding/json"
)

// String is an alias for sql.String data type
type String struct {
	Value string
	Valid bool
}

// Scan implements the Scanner interface for String
func (ns *String) Scan(value interface{}) error {
	switch value.(type) {
	case string:
		*ns = String{value.(string), true}
	case nil:
		*ns = String{"", false}
	default:
		return handleError(value)
	}

	return nil
}

// MarshalJSON for String
func (ns *String) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.Value)
}

// UnmarshalJSON for String
func (ns *String) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.Value)
	ns.Valid = (err == nil)
	return err
}
