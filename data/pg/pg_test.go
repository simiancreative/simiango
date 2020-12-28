package pg

import (
	"testing"

	m "github.com/simiancreative/simiango/mocks/data/pg"
)

var ConnMock *m.Conn

func init() {
	ConnMock = &m.Conn{}
	C = ConnMock
}

func TestQuery(t *testing.T) {
	query := "select * from widgets where id=$1"
	param := 42

	ConnMock.
		On("Query", Ctx, query, param).
		Return(nil, nil)

	C.Query(Ctx, query, param)

	ConnMock.AssertExpectations(t)
}
