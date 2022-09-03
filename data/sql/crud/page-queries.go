package crud

import "github.com/doug-martin/goqu/v9"

func (m *Model) query(filters Filters, orders Orders) *goqu.SelectDataset {
	ds := m.dialect.From(m.Table)
	ds = m.handleFilters(ds, filters)
	ds = m.handleOrder(ds, orders)
	ds = m.handleAugmentList(ds)

	return ds
}

func (m *Model) handleOrder(ds *goqu.SelectDataset, orders Orders) *goqu.SelectDataset {
	for key, value := range orders {
		ordr := goqu.C(key)

		if value == "asc" {
			ds = ds.Order(ordr.Asc())
		}

		if value == "dsc" {
			ds = ds.Order(ordr.Desc())
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
			ds = ds.Where(goqu.Ex(expressions.(Filters)))
		}

		if kind == OR {
			ds = ds.Where(goqu.ExOr(expressions.(Filters)))
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
