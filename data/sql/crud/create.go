package crud

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sanity-io/litter"
	mocks "github.com/simiancreative/simiango/mocks/data/sql"
)

func (m *Model) Create(params interface{}, destination interface{}) error {
	ds := m.dialect.Insert(m.Table).Rows(params)

	query, qParams, _ := ds.ToSQL()
	result, err := m.cx.Exec(query, qParams...)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	if err := m.One(destination, id); err != nil {
		return err
	}

	return nil
}

func (m *Model) SetupCreateTest(
	t *testing.T,
	params interface{},
) func() error {
	db, mock, clos := mocks.SetupAndMock("mysql", t)

	m.Initialize("mysql", db)

	ds := m.dialect.Insert(m.Table).Rows(params)
	query, _, _ := ds.ToSQL()

	mock.
		ExpectExec(regexp.QuoteMeta(query)).
		WillReturnResult(sqlmock.NewResult(456, 1))

	oneDs := m.oneQuery(456)
	oneQuery, _, _ := oneDs.ToSQL()

	litter.Dump(oneQuery)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(456)

	mock.
		ExpectQuery(regexp.QuoteMeta(oneQuery)).
		WillReturnRows(rows)

	m.cx = db

	return clos
}
