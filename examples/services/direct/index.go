package direct

import (
	"github.com/simiancreative/simiango/server"
	"github.com/simiancreative/simiango/service"

	m "github.com/simiancreative/simiango/data/mysql"
)

var Config = service.Config{
	Kind:      service.DIRECT,
	IsPrivate: true,
	Method:    "GET",
	Path:      "/direct",
	Direct:    direct,
}

type Product struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

func direct(req service.Req) (interface{}, error) {
	defer req.Timer.NewMetric("no_desc").Start().Stop()
	defer req.Timer.NewMetric("query").WithDesc("desc is optional").Start().Stop()

	rows := []Product{}
	err := m.Cx.Select(&rows, "select * from products")

	return rows, err
}

// dont forget to import your package in your main.go for initialization
// _ "path/to/project/stream"
func init() {
	server.AddService(Config)
}
