package main

import (
	"os"

	//
	// must be first, all env vars are loaded then others services that depend on
	// config can be loaded afterward
	//
	_ "github.com/simiancreative/simiango/config"
	"github.com/simiancreative/simiango/logger"

	"github.com/simiancreative/simiango/meta"

	_ "github.com/simiancreative/simiango/data/mssql"
	_ "github.com/simiancreative/simiango/data/mysql"
	_ "github.com/simiancreative/simiango/data/pg"

	"github.com/simiancreative/simiango/messaging/amqp"
	"github.com/simiancreative/simiango/messaging/kafka"
	"github.com/simiancreative/simiango/server"
	"github.com/simiancreative/simiango/stats/prometheus"

	_ "github.com/simiancreative/simiango/examples/services/assign"
	_ "github.com/simiancreative/simiango/examples/services/combinators"
	_ "github.com/simiancreative/simiango/examples/services/crud"
	_ "github.com/simiancreative/simiango/examples/services/crypt"
	_ "github.com/simiancreative/simiango/examples/services/direct"
	_ "github.com/simiancreative/simiango/examples/services/kafka/consume"
	_ "github.com/simiancreative/simiango/examples/services/kafka/consume-without-messages"
	_ "github.com/simiancreative/simiango/examples/services/kafka/ingest"
	_ "github.com/simiancreative/simiango/examples/services/mssql"
	_ "github.com/simiancreative/simiango/examples/services/mysql"
	_ "github.com/simiancreative/simiango/examples/services/mysql-page"
	_ "github.com/simiancreative/simiango/examples/services/param"
	_ "github.com/simiancreative/simiango/examples/services/pg"
	_ "github.com/simiancreative/simiango/examples/services/rabbit"
	_ "github.com/simiancreative/simiango/examples/services/sample"
	_ "github.com/simiancreative/simiango/examples/services/stream"
	_ "github.com/simiancreative/simiango/examples/services/unsafe"
)

func main() {
	logger.Printf("ENV STARTING AS: %v", os.Getenv("APP_ENV"))
	done, exit := meta.CatchSig()

	prometheus.Handle()

	server.EnableHealthCheck()
	server.SetCORS()
	server.AddPprof()

	go server.Start()

	_, startRabbit := os.LookupEnv("AMQP")
	if startRabbit {
		go amqp.Start()
	}

	_, startKafka := os.LookupEnv("KAFKA")
	if startKafka {
		go kafka.Start(done)
	}

	// keep main process open
	<-exit.Done()
}
