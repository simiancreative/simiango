package crud

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/doug-martin/goqu/v9"
	mocks "github.com/simiancreative/simiango/mocks/data/sql"
)

func (m *Model) One(destination interface{}, id interface{}) error {
	ds := m.oneQuery(id)
	query, params, _ := ds.ToSQL()
	if err := m.cx.Get(destination, query, params...); err != nil {
		return err
	}

	return nil
}

func (m *Model) handleAugmentOne(ds *goqu.SelectDataset) *goqu.SelectDataset {
	if m.AugmentOneQuery == nil {
		return ds
	}

	cb := *m.AugmentOneQuery
	return cb(ds)
}

func (m *Model) oneQuery(id interface{}) *goqu.SelectDataset {
	ds := m.query(Filters{}, Order{})
	ds = ds.Where(goqu.Ex{"id": id})
	ds = m.handleAugmentOne(ds)

	return ds
}

func (m *Model) SetupOneTest(
	t *testing.T,
) func() error {
	db, mock, clos := mocks.SetupAndMock("mysql", t)

	m.Initialize("mysql", db)

	ds := m.oneQuery(456)
	query, _, _ := ds.ToSQL()

	rows := sqlmock.NewRows([]string{"id"}).AddRow(456)

	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WillReturnRows(rows)

	m.cx = db

	return clos
}
