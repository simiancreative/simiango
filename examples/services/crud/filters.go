package crud

import (
	"strconv"

	"github.com/doug-martin/goqu/v9"
)

func asDescendants(ds *goqu.SelectDataset, values interface{}) *goqu.SelectDataset {
	return ds.Join(
		goqu.T("ancestors"),
		goqu.On(goqu.Ex{
			"ancestors.product_id": goqu.T("products").Col("id"),
			"ancestors.depth":      1,
		}),
	)
}

func ancestor(ds *goqu.SelectDataset, values interface{}) *goqu.SelectDataset {
	id, _ := strconv.Atoi(values.([]string)[0])
	return ds.Where(goqu.Ex{"ancestors.ancestor_id": id})
}

func name(ds *goqu.SelectDataset, values interface{}) *goqu.SelectDataset {
	name := values.([]string)[0]
	return ds.Where(goqu.ExOr{"name": name})
}
