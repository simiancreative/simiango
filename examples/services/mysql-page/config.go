package mysqlpage

import (
	"github.com/simiancreative/simiango/server"
	"github.com/simiancreative/simiango/service"
)

// HINT
// curl --location --request GET
// 'localhost:5100/mysql-example/products/page?page=1&size=25&pattern=Car&pattern2=Horse'

var Config = service.Config{
	Method: "GET",
	Path:   "/mysql-example/products/page",
	Build:  Build,
}

// dont forget to import your package in your main.go for initialization
// _ "path/to/project/stream"
func init() {
	server.AddService(Config)
}
