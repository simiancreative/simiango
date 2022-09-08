package crud

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
)

func (m *Model) query(filters Filters, order Order) *goqu.SelectDataset {
	ds := m.dialect.From(m.Table)
	ds = m.handleSelect(ds)
	ds = m.handleFilters(ds, filters)
	ds = m.handleOrder(ds, order)

	return ds
}

func (m *Model) handleSelect(ds *goqu.SelectDataset) *goqu.SelectDataset {
	if len(m.Columns) == 0 {
		return ds.Select(fmt.Sprintf("%v.%v", m.Table, "*"))
	}

	return ds.Select(m.Columns...)
}

func (m *Model) handleOrder(ds *goqu.SelectDataset, order Order) *goqu.SelectDataset {
	for _, col := range order {
		ordr := goqu.C(col.Name)

		if col.Direction == ASC {
			ds = ds.OrderAppend(ordr.Asc())
		}

		if col.Direction == DSC {
			ds = ds.OrderAppend(ordr.Desc())
		}
	}

	return ds
}

func (m *Model) handleFilters(ds *goqu.SelectDataset, filters Filters) *goqu.SelectDataset {
	for key, value := range filters {
		filter, ok := m.Expressions[key]
		if !ok {
			continue
		}

		ds = filter(ds, value)
	}

	return ds
}

func (m *Model) handleAugmentList(ds *goqu.SelectDataset) *goqu.SelectDataset {
	if m.AugmentListQuery == nil {
		return ds
	}

	cb := *m.AugmentListQuery
	return cb(ds)
}

func (m *Model) pageQueries(ds *goqu.SelectDataset, page int, size int) (*goqu.SelectDataset, *goqu.SelectDataset) {
	countQuery := ds.Select(goqu.COUNT("*"))
	ds = ds.Limit(uint(size)).Offset(uint((page - 1) * size))

	return countQuery, ds
}
