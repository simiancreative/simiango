package nulls

import "fmt"

// BitBool is an implementation of a bool for the MySQL type BIT(1).
// This type allows you to avoid wasting an entire byte for MySQL's boolean type TINYINT.
type BitBool struct {
	Value bool
	Valid bool
}

// Scan implements the sql.Scanner interface,
// and turns the bitfield incoming from MySQL into a BitBool
func (b *BitBool) Scan(value interface{}) error {
	switch value.(type) {
	case []byte:
		*b = BitBool{value.([]byte)[0] == 1, true}
	case bool:
		*b = BitBool{value.(bool), true}
	case nil:
		*b = BitBool{false, false}
	default:
		return handleError(value)
	}

	return nil
}

// MarshalJSON for BitBool
func (b *BitBool) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return []byte("null"), nil
	}

	val := fmt.Sprintf("%t", b.Value)
	return []byte(val), nil
}

// UnmarshalJSON for Bool
func (b *BitBool) UnmarshalJSON(v []byte) error {
	s := string(v)

	if s != "true" && s != "false" {
		b.Valid = false
		return nil
	}

	b.Value = s == "true"
	b.Valid = true

	return nil
}
