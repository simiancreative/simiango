package requestreceiver

import (
	"github.com/simiancreative/simiango/server"
	"github.com/simiancreative/simiango/service"
)

var Config = service.Config{
	Kind:   service.DIRECT,
	Method: "POST",
	Path:   "/request-receiver",
	Direct: direct,
	Input:  input,
}

type Product struct {
	ID   string `json:"id"   db:"id"`
	Name string `json:"name" db:"name"`
}

func direct(req service.Req) (interface{}, error) {
	return req.Input, nil
}

func input() interface{} {
	return &Product{}
}

func init() {
	server.AddService(Config)
}
