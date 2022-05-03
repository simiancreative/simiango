package mocks

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	sqlx "github.com/jmoiron/sqlx"
)

func SetupAndMock(driver string, t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func() error) {
	db, mock, err := sqlmock.New()

	ConnXMock := sqlx.NewDb(db, driver)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return ConnXMock, mock, ConnXMock.Close
}
