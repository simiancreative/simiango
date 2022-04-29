package main

import (

	//
	// must be first, all env vars are loaded then others services that depend on
	// config can be loaded afterward
	//
	"fmt"
	"os"

	_ "github.com/simiancreative/simiango/config"
	"github.com/simiancreative/simiango/server"
	"github.com/simiancreative/simiango/stats/prometheus"

	_ "github.com/simiancreative/simiango/data/mssql"
	_ "github.com/simiancreative/simiango/data/mysql"
	_ "github.com/simiancreative/simiango/data/pg"

	"github.com/simiancreative/simiango/messaging/amqp"
	"github.com/simiancreative/simiango/messaging/kafka"

	_ "github.com/simiancreative/simiango/examples/services/assign"
	_ "github.com/simiancreative/simiango/examples/services/crypt"
	_ "github.com/simiancreative/simiango/examples/services/direct"
	_ "github.com/simiancreative/simiango/examples/services/kafka"
	_ "github.com/simiancreative/simiango/examples/services/mssql"
	_ "github.com/simiancreative/simiango/examples/services/mysql"
	_ "github.com/simiancreative/simiango/examples/services/mysql-page"
	_ "github.com/simiancreative/simiango/examples/services/param"
	_ "github.com/simiancreative/simiango/examples/services/pg"
	_ "github.com/simiancreative/simiango/examples/services/rabbit"
	_ "github.com/simiancreative/simiango/examples/services/sample"
	_ "github.com/simiancreative/simiango/examples/services/stream"
)

func main() {
	prometheus.Handle()

	server.EnableHealthCheck()
	server.SetCORS()

	go server.Start()

	_, startRabbit := os.LookupEnv("AMQP")
	if startRabbit {
		go amqp.Start()
	}

	_, startKafka := os.LookupEnv("KAFKA")
	if startKafka {
		//go kafka.Start()
		done := kafka.Start()
		fmt.Println(<-done)
		fmt.Println("After done in example main")
	}

	// keep main process open
	//select {}
}
