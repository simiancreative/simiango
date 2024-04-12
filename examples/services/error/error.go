package error

import (
	"fmt"

	"github.com/simiancreative/simiango/server"
	"github.com/simiancreative/simiango/service"
)

var Config = service.Config{
	Kind:   service.DIRECT,
	Method: "GET",
	Path:   "/error",
	Direct: direct,
}

func direct(req service.Req) (interface{}, error) {
	return nil, fmt.Errorf("this is an error")
}

// dont forget to import your package in your main.go for initialization
// _ "path/to/project/direct"
func init() {
	server.AddService(Config)
}
