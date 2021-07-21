package mongoexample

import (
	"github.com/simiancreative/simiango/service"
)

var Config = service.Config{
	Method: "GET",
	Path:   "/mongo-example/products",
	Build:  Build,
}
