package mysqlexample

import (
	"github.com/simiancreative/simiango/service"
)

var Config = service.Config{
	Method: "GET",
	Path:   "/mysql-example/products",
	Build:  Build,
}
