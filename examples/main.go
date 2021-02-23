package main

import (
	_ "github.com/simiancreative/simiango/config"
	_ "github.com/simiancreative/simiango/examples/docs"

	_ "github.com/simiancreative/simiango/data/mssql"
	_ "github.com/simiancreative/simiango/data/pg"

	"github.com/simiancreative/simiango/server"
	"github.com/simiancreative/simiango/service"

	mssql "github.com/simiancreative/simiango/examples/services/mssql"
	pg "github.com/simiancreative/simiango/examples/services/pg"
	sample "github.com/simiancreative/simiango/examples/services/sample"
)

// @title Sample API
// @version 1.0
// @description This is a sample service for simian go
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /

func main() {
	server.EnableHealthCheck()
	server.EnableSwagger()
	server.SetCORS()

	server.AddServices(service.Collection{
		mssql.Config,
		pg.Config,
		sample.Config,
	})

	server.Start()
}
