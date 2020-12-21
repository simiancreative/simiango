package main

import (
	_ "github.com/simiancreative/simiango/config"
	_ "github.com/simiancreative/simiango/data/redis"
	_ "github.com/simiancreative/simiango/sample/docs"

	"github.com/simiancreative/simiango/server"
	"github.com/simiancreative/simiango/service"

	sample "github.com/simiancreative/simiango/sample/sample"
)

// @title Sample API
// @version 1.0
// @description This is a sample service for simian go
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /

func main() {
	server.EnableSwagger()
	server.SetCORS()

	server.AddServices(service.Collection{
		sample.Config,
	})

	server.Start()
}
