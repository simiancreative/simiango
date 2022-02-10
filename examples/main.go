package main

import (
	"os"

	_ "github.com/simiancreative/simiango/config"

	_ "github.com/simiancreative/simiango/data/mssql"
	_ "github.com/simiancreative/simiango/data/mysql"
	_ "github.com/simiancreative/simiango/data/pg"

	"github.com/simiancreative/simiango/messaging/amqp"
	"github.com/simiancreative/simiango/messaging/kafka"
	"github.com/simiancreative/simiango/server"
	"github.com/simiancreative/simiango/service"

	_ "github.com/simiancreative/simiango/examples/services/crypt"
	mssql "github.com/simiancreative/simiango/examples/services/mssql"
	mysql "github.com/simiancreative/simiango/examples/services/mysql"
	pg "github.com/simiancreative/simiango/examples/services/pg"
	sample "github.com/simiancreative/simiango/examples/services/sample"
	_ "github.com/simiancreative/simiango/examples/services/stream"

	_ "github.com/simiancreative/simiango/examples/services/kafka"
	_ "github.com/simiancreative/simiango/examples/services/param"
	_ "github.com/simiancreative/simiango/examples/services/rabbit"
	_ "github.com/simiancreative/simiango/examples/services/rabbit2"
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

	_, startRabbit := os.LookupEnv("AMQP")
	if startRabbit {
		go amqp.Start()
	}

	_, startKafka := os.LookupEnv("KAFKA")
	if startKafka {
		go kafka.Start()
	}

	// keep main process open
	select {}
}
