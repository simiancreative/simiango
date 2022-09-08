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
	ds = m.handleAugmentList(ds)

	return ds
}

func (m *Model) handleSelect(ds *goqu.SelectDataset) *goqu.SelectDataset {
	if len(m.Columns) == 0 {
		return ds.Select(fmt.Sprintf("%v.%v", m.Table, "*"))
	}

	return ds.Select(m.Columns...)
}

func (m *Model) handleOrder(ds *goqu.SelectDataset, order Order) *goqu.SelectDataset {
	for key, value := range order {
		ordr := goqu.C(key)

		if value == "asc" {
			ds = ds.OrderAppend(ordr.Asc())
		}

		if value == "dsc" {
			ds = ds.OrderAppend(ordr.Desc())
		}
	}

	return ds
}

func (m *Model) handleFilters(ds *goqu.SelectDataset, filters Filters) *goqu.SelectDataset {
	for key, value := range filters {
		filter, ok := m.Filterable[key]
		if !ok {
			continue
		}

		kind, expressions := filter(value)
		if kind == AND {
			ds = ds.Where(goqu.Ex(expressions.(Exp)))
		}

		if kind == OR {
			ds = ds.Where(goqu.ExOr(expressions.(Exp)))
		}

		if kind == JOIN {
			join := expressions.(Joinable)

			ds = ds.Join(
				goqu.T(join.Table),
				goqu.On(goqu.Ex(join.Exp)),
			)
		}
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
