package main

import (
	"github.com/simiancreative/simiango/server"
	"github.com/simiancreative/simiango/service"

	sample "github.com/simiancreative/simiango/sample/sample"
)

func main() {
	server.SetCORS()

	server.AddServices(service.Collection{
		sample.Config,
	})

	server.Start()
}
