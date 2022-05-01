package unsafe

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	sqlx "github.com/jmoiron/sqlx"

	"github.com/stretchr/testify/assert"
)

var ConnXMock *sqlx.DB

var mockItem = Item{"id": int64(23), "name": "Floppy diskette"}

var mockContent = Content{
	// demonstrating that this is not a type safe solution and has the potential
	// to cause runtime errors
	Item{"id": int64(23), "name": "Floppy diskette"},
	Item{"id": "one", "name": int64(42)},
}

func setup(t *testing.T) func() error {
	db, mock, err := sqlmock.New()

	ConnXMock = sqlx.NewDb(db, "mysql")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(int64(23), "Floppy diskette").
		AddRow("one", int64(42))

	mock.
		ExpectQuery("SELECT id FROM devices WHERE gateway_id = ?").
		WithArgs("1234").
		WillReturnRows(rows)

	return ConnXMock.Close
}

func TestUnsafeSelect(t *testing.T) {
	defer setup(t)()

	u := Unsafe{Cx: ConnXMock}

	content, err := u.UnsafeSelect(
		"SELECT id FROM devices WHERE gateway_id = ?",
		"1234",
	)

	assert.NoError(t, err, "has error")
	assert.Equal(t, content, mockContent, "content is different")
}

func TestUnsafeGet(t *testing.T) {
	defer setup(t)()

	u := Unsafe{Cx: ConnXMock}

	content, err := u.UnsafeGet(
		"SELECT id FROM devices WHERE gateway_id = ?",
		"1234",
	)

	assert.NoError(t, err)
	assert.Equal(t, content, mockItem)
}
