package crud

import "github.com/doug-martin/goqu/v9"

func (m *Model) One(destination interface{}, id interface{}) error {
	ds := m.query(Filters{}, Orders{})
	ds = ds.Where(goqu.Ex{"id": id})
	ds = m.handleAugmentOne(ds)

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
