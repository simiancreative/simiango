package pg_test

import (
	"testing"

	pg "github.com/simiancreative/simiango/data/pg"
	m "github.com/simiancreative/simiango/mocks/data/sql"
)

var ConnXMock *m.ConnX

func init() {
	ConnXMock = &m.ConnX{}
	pg.Cx = ConnXMock
}

func TestQueryX(t *testing.T) {
	query := "select * from widgets where id=$1"
	param := 42

	ConnXMock.
		On("Queryx", query, param).
		Return(nil, nil)

	pg.Cx.Queryx(query, param)

	ConnXMock.AssertExpectations(t)
}
