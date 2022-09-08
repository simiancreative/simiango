package crud

import (
	"strconv"

	"github.com/doug-martin/goqu/v9"
	"github.com/simiancreative/simiango/data/sql/crud"
)

func asDescendants(values interface{}) (crud.Kind, interface{}) {
	return crud.JOIN, crud.Joinable{
		Table: "ancestors",
		Exp: crud.Exp{
			"ancestors.product_id": goqu.T("products").Col("id"),
			"ancestors.depth":      1,
		},
	}
}

func ancestor(values interface{}) (crud.Kind, interface{}) {
	id, _ := strconv.Atoi(values.([]string)[0])
	return crud.AND, crud.Exp{"ancestors.ancestor_id": id}
}

func name(values interface{}) (crud.Kind, interface{}) {
	name := values.([]string)[0]

	return crud.OR, crud.Exp{"name": name}
}
