package pgservice

import (
	"github.com/simiancreative/simiango/service"
)

var Config = service.Config{
	Method: "GET",
	Path:   "/pg-example/products",
	Build:  Build,
}
