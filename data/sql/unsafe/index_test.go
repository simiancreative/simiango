package unsafe

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	sqlx "github.com/jmoiron/sqlx"
	"github.com/simiancreative/simiango/mocks/data/sql"

	"github.com/stretchr/testify/assert"
)

var ConnXMock *sqlx.DB

var mockContent = []interface{}{
	// demonstrating that this is not a type safe solution and has the potential
	// to cause runtime errors
	Item{"id": int64(42), "name": "Floppy diskette"},
	Item{"id": int64(12), "name": "SSD"},
}

func setupQuery(t *testing.T) func() error {
	db, mock, clos := mocks.SetupAndMock("mysql", t)
	ConnXMock = db

	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, i := range mockContent {
		rows.AddRow(i.(Item)["id"], i.(Item)["name"])
	}

	mock.
		ExpectQuery("SELECT id, name FROM devices WHERE gateway_id = ?").
		WithArgs("1234").
		WillReturnRows(rows)

	return clos
}

func TestUnsafeSelect(t *testing.T) {
	defer setupQuery(t)()

	u := Unsafe{Cx: ConnXMock}

	result, err := u.UnsafeSelect(
		"SELECT id, name FROM devices WHERE gateway_id = ?",
		"1234",
	)

	assert.NoError(t, err, "has error")
	assert.Equal(t, result.Content, mockContent, "content is different")
}

func TestUnsafeGet(t *testing.T) {
	defer setupQuery(t)()

	u := Unsafe{Cx: ConnXMock}

	result, err := u.UnsafeGet(
		"SELECT id, name FROM devices WHERE gateway_id = ?",
		"1234",
	)

	mockResult := Result{
		Columns: []Column{{Name: "id"}, {Name: "name"}},
		Content: []interface{}{mockContent[0]},
	}

	assert.NoError(t, err)
	assert.Equal(t, mockResult, result)
}
