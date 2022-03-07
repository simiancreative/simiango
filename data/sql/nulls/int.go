package nulls

import (
	"encoding/json"
)

type Int64 struct {
	Value int64
	Valid bool
}

func (nt *Int64) Scan(value interface{}) error {
	switch value.(type) {
	case int64:
		*nt = Int64{value.(int64), true}
	case nil:
		*nt = Int64{int64(0), false}
	default:
		return handleError(value)
	}

	return nil
}

func (nt *Int64) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Value)
}

func (nt *Int64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &nt.Value)
	nt.Valid = (err == nil)
	return err
}
