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

	err := u.Cx.Select(&items, query, params...)

	return items, err
}

// UnsafeGet is used when expecting a single instance of output values where
// types are unknown
func (u *Unsafe) UnsafeGet(query string, params ...interface{}) (Item, error) {
	item := Item{}

	err := u.Cx.Get(&item, query, params...)

	return item, err
}
