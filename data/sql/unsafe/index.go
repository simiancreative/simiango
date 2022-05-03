package unsafe

import (
	"github.com/simiancreative/simiango/data/sql"
)

// Unsafe sets up a sql connection interface for later unsafe queries
type Unsafe struct {
	Cx sql.ConnX
}

// UnsafeSelect is used when expecting multiple instances of output values where
// types are unknown
func (u *Unsafe) UnsafeSelect(query string, params ...interface{}) (Result, error) {
	result := Result{}

	rows, err := u.Cx.Queryx(query, params...)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		if err := result.addColumns(rows); err != nil {
			return result, err
		}

		if err := result.addItem(rows); err != nil {
			return result, err
		}
	}

	rows.Close()

	return result, err
}

// UnsafeGet is used when expecting a single instance of output values where
// types are unknown
func (u *Unsafe) UnsafeGet(query string, params ...interface{}) (Result, error) {
	result := Result{}

	rows, err := u.Cx.Queryx(query, params...)
	if err != nil {
		return result, err
	}

	rows.Next()

	if err := result.addColumns(rows); err != nil {
		return result, err
	}

	if err := result.addItem(rows); err != nil {
		return result, err
	}

	rows.Close()

	return result, err
}
