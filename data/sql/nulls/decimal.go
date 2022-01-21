package nulls

import (
	"encoding/json"
)

type Decimal struct {
	Value float64
	Valid bool
}

func (nt *Decimal) Scan(value interface{}) error {
	switch value.(type) {
	case []uint8:
		*nt = Decimal{Float64fromUint8s(value.([]uint8)), true}
	case float64:
		*nt = Decimal{value.(float64), true}
	default:
		*nt = Decimal{float64(0), false}
	}

	if !nt.Valid {
		return handleError(value)
	}

	return nil
}

// MarshalJSON for Decimal
func (nt *Decimal) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Value)
}

// UnmarshalJSON for Decimal
func (nt *Decimal) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &nt.Value)
	nt.Valid = (err == nil)
	return err
}
