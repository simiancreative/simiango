package sample

import (
	"github.com/simiancreative/simiango/service"
)

var Config = service.Config{
	Method: "POST",
	Path:   "/sample/:id",
	Build:  Build,
}
