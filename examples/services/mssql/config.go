package mssqlexample

import (
	"github.com/simiancreative/simiango/service"
)

var Config = service.Config{
	Method: "GET",
	Path:   "/mssql-example/products",
	Build:  Build,
}
