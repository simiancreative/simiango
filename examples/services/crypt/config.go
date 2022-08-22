package crypt

import (
	"os"

	"github.com/simiancreative/simiango/cryptkeeper"
	"github.com/simiancreative/simiango/server"
	"github.com/simiancreative/simiango/service"
)

var Config = service.Config{
	IsPrivate: false,
	Method:    "GET",
	Path:      "/crypt",
	Build:     Build,
}

// dont forget to import your package in your main.go for initialization
// _ "path/to/project/crypt"
func init() {
	cryptkeeper.SetKey(os.Getenv("CRYPTED_TOKEN"))
	server.AddService(Config)
}
