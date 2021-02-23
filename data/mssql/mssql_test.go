package mssql

import (
	"context"
	"testing"

	m "github.com/simiancreative/simiango/mocks/data/mssql"
)

var ConnXMock *m.ConnX

func init() {
	Ctx = context.Background()

	ConnXMock = &m.ConnX{}
	Cx = ConnXMock

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
