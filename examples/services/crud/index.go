package crud

import (
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

	Filterable: crud.Filterable{
		"as_descendants": asDescendants,
		"ancestor":       ancestor,
		"name":           name,
	},
}

// dont forget to import your package in your main.go for initialization
// _ "path/to/project/direct"
func init() {
	c.Initialize(mysql.Cx, "mysql")
	server.AddService(Config)
}
