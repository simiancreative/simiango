package sample

import (
	"github.com/simiancreative/simiango/server"
	"github.com/simiancreative/simiango/service"
)

var Config = service.Config{
	IsPrivate: true,
	Method:    "POST",
	Path:      "/sample/:id",
	Build:     Build,
}

// dont forget to import your package in your main.go for initialization
// _ "path/to/project/stream"
func init() {
	server.AddService(Config)
}
