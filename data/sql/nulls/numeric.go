package nulls

import (
	"encoding/json"
)

type Numeric struct {
	Value float64
	Valid bool
}

func (nt *Numeric) Scan(value interface{}) error {
	switch value.(type) {
	case []uint8:
		*nt = Numeric{Float64fromUint8s(value.([]uint8)), true}
	case float64:
		*nt = Numeric{value.(float64), true}
	case nil:
		*nt = Numeric{float64(0), false}
	default:
		return handleError(value)
	}

	return nil
}

// MarshalJSON for Numeric
func (nt *Numeric) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Value)
}

// UnmarshalJSON for Numeric
func (nt *Numeric) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &nt.Value)
	nt.Valid = (err == nil)
	return err
}
