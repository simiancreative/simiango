package pg

import (
	"testing"

	m "github.com/simiancreative/simiango/mocks/data/pg"
)

var ConnMock *m.Conn
var ConnXMock *m.ConnX

func init() {
	ConnMock = &m.Conn{}

	ConnXMock = &m.ConnX{}
	Cx = ConnXMock
}

func TestQueryX(t *testing.T) {
	query := "select * from widgets where id=$1"
	param := 42

	ConnXMock.
		On("Queryx", query, param).
		Return(nil, nil)

	Cx.Queryx(query, param)

	ConnXMock.AssertExpectations(t)
}
