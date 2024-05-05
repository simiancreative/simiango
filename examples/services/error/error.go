package error

import (
	"fmt"

	"github.com/jackc/pgtype"

	"github.com/simiancreative/simiango/server"
	"github.com/simiancreative/simiango/service"
)

// JUST AN ERROR

var Config = service.Config{
	Kind:   service.DIRECT,
	Method: "GET",
	Path:   "/error",
	Direct: direct,
}

func direct(req service.Req) (interface{}, error) {
	return nil, fmt.Errorf("this is an error")
}

// BAD JSON

var badJsonConfig = service.Config{
	Kind:   service.DIRECT,
	Method: "GET",
	Path:   "/error/bad-json",
	Direct: badJson,
}

type Product struct {
	ID   string       `json:"id"   db:"id"`
	Meta pgtype.JSONB `json:"meta" db:"meta"`
}

func badJson(_ service.Req) (interface{}, error) {
	return Product{}, nil
}

// dont forget to import your package in your main.go for initialization
// _ "path/to/project/direct"
func init() {
	server.AddServices([]service.Config{Config, badJsonConfig})
}
