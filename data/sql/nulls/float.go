package nulls

import (
	"encoding/json"
)

type Float64 struct {
	Value float64
	Valid bool
}

func (nt *Float64) Scan(value interface{}) error {
	switch value.(type) {
	case float64:
		*nt = Float64{value.(float64), true}
	case nil:
		*nt = Float64{float64(0), false}
	default:
		return handleError(value)
	}

	return nil
}

func (nt *Float64) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Value)
}

func (nt *Float64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &nt.Value)
	nt.Valid = (err == nil)
	return err
}
