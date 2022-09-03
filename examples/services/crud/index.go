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
	Path:   "/model",
	Direct: direct,
}

var c = crud.Model{
	Table: "products",
	Filterable: crud.Filterable{
		"name": func(value interface{}) (crud.WhereKind, interface{}) {
			return crud.AND, crud.Filters{"name": value}
		},
	},
}

func direct(req service.Req) (interface{}, error) {
	result := []interface{}{}
	items := &Products{}

	page, err := c.Page(items, crud.Filters{"name": "Cars"}, crud.Orders{}, 1, 5)
	result = append(result, page)

	item := &Product{}
	err = c.One(item, 6)
	result = append(result, *item)

	params := &ProductProperties{Name: "Roller Blades"}
	err = c.Create(params, item)
	result = append(result, *item)

	params = &ProductProperties{Name: "Cars"}
	err = c.Update(item.ID, params, item)
	result = append(result, *item)

	err = c.Delete(item.ID)

	return result, err
}

// dont forget to import your package in your main.go for initialization
// _ "path/to/project/direct"
func init() {
	c.Initialize(mysql.Cx, "mysql")
	server.AddService(Config)
}
