package combinators

import (
	"github.com/simiancreative/simiango/server"
	"github.com/simiancreative/simiango/service"

	"github.com/simiancreative/simiango/data/pg"
	"github.com/simiancreative/simiango/data/sql/combinators"
)

type Product struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type Products []Product

var Config = service.Config{
	Kind:   service.DIRECT,
	Method: "GET",
	Path:   "/combinators/named-get",
	Direct: direct,
}

func direct(_ service.Req) (interface{}, error) {
	products := Products{}

	model := combinators.New(pg.Cx)
	err := model.NamedSelect(
		&products,
		"SELECT * FROM products where name = :name",
		Product{Name: "Truck"},
	)

	return products, err
}

// dont forget to import your package in your main.go for initialization
// _ "path/to/project/direct"
func init() {
	server.AddService(Config)
}
