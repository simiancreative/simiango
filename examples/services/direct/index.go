package direct

import (
	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/server"
	"github.com/simiancreative/simiango/service"

	m "github.com/simiancreative/simiango/data/mysql"
)

var Config = service.Config{
	Kind:   service.DIRECT,
	Method: "GET",
	Path:   "/direct",
	Direct: direct,
	After:  after,
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

// after is run asynchronously after so as to not delay the service return
func after(config service.Config, req service.Req) {
	logger.Printf("in after: (%v %v) -> (%v) %v",
		config.Method, config.Path, "direct", req.Timer.String(),
	)
}

// dont forget to import your package in your main.go for initialization
// _ "path/to/project/direct"
func init() {
	server.AddService(Config)
}
