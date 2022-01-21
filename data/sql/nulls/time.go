package nulls

import (
	"fmt"
	"time"
)

type Time struct {
	time.Time
	Valid bool
}

func (nt *Time) Scan(value interface{}) error {
	switch value.(type) {
	case time.Time:
		*nt = Time{value.(time.Time), true}
	default:
		*nt = Time{time.Time{}, false}
	}

	if !nt.Valid {
		return handleError(value)
	}

	return nil
}

// MarshalJSON for Time
func (nt *Time) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	val := fmt.Sprintf("\"%s\"", nt.Time.Format(time.RFC3339))
	return []byte(val), nil
}

// UnmarshalJSON for Time
func (nt *Time) UnmarshalJSON(b []byte) error {
	s := string(b)
	// s = Stripchars(s, "\"")

	x, err := time.Parse(time.RFC3339, s)
	if err != nil {
		nt.Valid = false
		return err
	}

	nt.Time = x
	nt.Valid = true
	return nil
}
