package sample

import (
	"github.com/simiancreative/simiango/service"
)

var Config = service.Config{
	IsPrivate: true,
	Method:    "POST",
	Path:      "/sample/:id",
	Build:     Build,
}
