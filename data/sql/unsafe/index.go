package unsafe

import (
	"github.com/simiancreative/simiango/data/sql"
)

// Unsafe sets up a sql connection interface for later unsafe queries
type Unsafe struct {
	Cx sql.ConnX
}

// Item is used as a container for sql rows with unknown values
type Item map[string]interface{}

// Content is used as a container for Unsafe Items
type Content []Item

// UnsafeSelect is used when expecting multiple instances of output values where
// types are unknown
func (u *Unsafe) UnsafeSelect(query string, params ...interface{}) (Content, error) {
	items := Content{}

	rows, err := u.Cx.Queryx(query, params...)
	if err != nil {
		return items, err
	}

	for rows.Next() {
		item := Item{}
		err := rows.MapScan(item)
		if err != nil {
			return items, err
		}

		items = append(items, item)
	}

	rows.Close()

	return items, err
}

// UnsafeGet is used when expecting a single instance of output values where
// types are unknown
func (u *Unsafe) UnsafeGet(query string, params ...interface{}) (Item, error) {
	item := Item{}

	rows, err := u.Cx.Queryx(query, params...)
	if err != nil {
		return item, err
	}

	rows.Next()

	err = rows.MapScan(item)
	if err != nil {
		return item, err
	}

	rows.Close()

	return item, err
}
