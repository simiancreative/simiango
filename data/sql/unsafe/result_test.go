package unsafe

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/simiancreative/simiango/mocks/data/sql"
)

func setupBinary(t *testing.T) func() error {
	db, mock, clos := mocks.SetupAndMock("mysql", t)
	ConnXMock = db

	byteRows := sqlmock.NewRows([]string{"id", "name"})
	byteRows.AddRow(1, []byte("a binary value"))

	mock.
		ExpectQuery("SELECT cost FROM products").
		WillReturnRows(byteRows)

	return clos
}

func TestUnsafeResultaddItem(t *testing.T) {
	defer setupBinary(t)()

	u := Unsafe{Cx: ConnXMock}

	result, err := u.UnsafeGet("SELECT cost FROM products")

	assert.NoError(t, err)
	assert.Equal(t, "a binary value", result.Content[0].(Item)["name"])
}
