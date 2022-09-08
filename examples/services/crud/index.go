package crud

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/simiancreative/simiango/server"
	"github.com/simiancreative/simiango/service"

	"github.com/simiancreative/simiango/data/mysql"
	"github.com/simiancreative/simiango/data/sql/crud"
)

var Config = service.Config{
	Kind:   service.DIRECT,
	Method: "GET",
	Path:   "/crud",
	Direct: result,
}

var c = crud.Model{
	Table: "products",
	Columns: crud.Columns{
		goqu.T("ancestors").Col("depth"),
		goqu.T("products").Col("*"),
	},

	Expressions: crud.Expressions{
		"as_descendants": asDescendants,
		"ancestor":       ancestor,
		"name":           name,
	},
}

// dont forget to import your package in your main.go for initialization
// _ "path/to/project/direct"
func init() {
	c.Initialize("mysql", mysql.Cx)
	server.AddService(Config)
}
