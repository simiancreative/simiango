package sql

import (
	"encoding/json"
	"reflect"
)

type NullFloat64 struct {
	Value float64
	Valid bool
}

func (nt *NullFloat64) Scan(value interface{}) error {
	if reflect.TypeOf(value) == nil {
		*nt = NullFloat64{0, false}
	} else {
		*nt = NullFloat64{value.(float64), true}
	}

	return nil
}

func (nt *NullFloat64) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Value)
}

func (nt *NullFloat64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &nt.Value)
	nt.Valid = (err == nil)
	return err
}
