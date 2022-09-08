package crud

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/doug-martin/goqu/v9"
	"github.com/sanity-io/litter"
	mocks "github.com/simiancreative/simiango/mocks/data/sql"
)

func (m *Model) Update(id interface{}, params interface{}, destination interface{}) error {
	ds := m.updateQuery(id, params)
	query, qParams, _ := ds.ToSQL()
	if _, err := m.cx.Exec(query, qParams...); err != nil {
		return err
	}

	if err := m.One(destination, id); err != nil {
		return err
	}

	return nil
}

func (m *Model) updateQuery(id interface{}, params interface{}) *goqu.UpdateDataset {
	return m.dialect.Update(m.Table).Set(params).Where(goqu.Ex{"id": id})
}

func (m *Model) SetupUpdateTest(
	t *testing.T,
	params interface{},
	destination interface{},
) func() error {
	db, mock, clos := mocks.SetupAndMock("mysql", t)

	m.Initialize("mysql", db)

	ds := m.updateQuery(456, params)
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
