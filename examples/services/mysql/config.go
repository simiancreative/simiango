package mysqlexample

import (
	"github.com/simiancreative/simiango/server"
	"github.com/simiancreative/simiango/service"
)

var Config = service.Config{
	Method: "GET",
	Path:   "/mysql-example/products",
	Build:  Build,
}

// dont forget to import your package in your main.go for initialization
// _ "path/to/project/stream"
func init() {
	server.AddService(Config)
}
