package sample

import (
	"simian/service"
)

var Config = service.Config{
	Method: "POST",
	Path:   "/sample/:id",
	Build:  Build,
}
