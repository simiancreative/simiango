package crud

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/doug-martin/goqu/v9"
	mocks "github.com/simiancreative/simiango/mocks/data/sql"
)

func (m *Model) Delete(id interface{}) error {
	ds := m.deleteQuery(id)
	query, qParams, _ := ds.ToSQL()
	_, err := m.cx.Exec(query, qParams...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) deleteQuery(id interface{}) *goqu.DeleteDataset {
	return m.dialect.Delete(m.Table).Where(goqu.Ex{"id": id})
}

func (m *Model) SetupDeleteTest(
	t *testing.T,
) func() error {
	db, mock, clos := mocks.SetupAndMock("mysql", t)

	m.Initialize("mysql", db)

	ds := m.deleteQuery(456)
	query, _, _ := ds.ToSQL()

	mock.
		ExpectExec(regexp.QuoteMeta(query)).
		WillReturnResult(sqlmock.NewResult(456, 1))

	m.cx = db

	return clos
}
