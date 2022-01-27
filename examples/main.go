package main

import (
	_ "github.com/simiancreative/simiango/config"

	_ "github.com/simiancreative/simiango/data/mssql"
	_ "github.com/simiancreative/simiango/data/mysql"
	_ "github.com/simiancreative/simiango/data/pg"

	"github.com/simiancreative/simiango/messaging/amqp"
	"github.com/simiancreative/simiango/server"
	"github.com/simiancreative/simiango/service"

	_ "github.com/simiancreative/simiango/examples/services/crypt"
	mssql "github.com/simiancreative/simiango/examples/services/mssql"
	mysql "github.com/simiancreative/simiango/examples/services/mysql"
	pg "github.com/simiancreative/simiango/examples/services/pg"
	sample "github.com/simiancreative/simiango/examples/services/sample"
	_ "github.com/simiancreative/simiango/examples/services/stream"

	_ "github.com/simiancreative/simiango/examples/services/rabbit"
)

func main() {
	server.EnableHealthCheck()
	server.SetCORS()

	server.AddServices(service.Collection{
		mssql.Config,
		pg.Config,
		mysql.Config,
		sample.Config,
	})

	go server.Start()
	go amqp.Start()

	// keep main process open
	select {}
}
