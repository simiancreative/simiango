package sql

import (
	"encoding/json"
	"reflect"
)

type NullInt64 struct {
	Value int64
	Valid bool
}

func (nt *NullInt64) Scan(value interface{}) error {
	if reflect.TypeOf(value) == nil {
		*nt = NullInt64{0, false}
	} else {
		*nt = NullInt64{value.(int64), true}
	}

	return nil
}

func (nt *NullInt64) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Value)
}

func (nt *NullInt64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &nt.Value)
	nt.Valid = (err == nil)
	return err
}
