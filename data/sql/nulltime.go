package sql

import (
	"fmt"
	"reflect"
	"time"
)

type NullTime struct {
	time.Time
	Valid bool
}

func (nt *NullTime) Scan(value interface{}) error {
	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*nt = NullTime{nt.Time, false}
	} else {
		*nt = NullTime{value.(time.Time), true}
	}

	return nil
}

// MarshalJSON for NullTime
func (nt *NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	val := fmt.Sprintf("\"%s\"", nt.Time.Format(time.RFC3339))
	return []byte(val), nil
}

// UnmarshalJSON for NullTime
func (nt *NullTime) UnmarshalJSON(b []byte) error {
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
