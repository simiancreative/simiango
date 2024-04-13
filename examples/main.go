package main

import (
	"os"

	"github.com/simiancreative/simiango/config"
	"github.com/simiancreative/simiango/data/mssql"
	"github.com/simiancreative/simiango/data/mysql"
	"github.com/simiancreative/simiango/data/pg"
	"github.com/simiancreative/simiango/examples/lib"
	_ "github.com/simiancreative/simiango/examples/services/assign"
	_ "github.com/simiancreative/simiango/examples/services/combinators"
	_ "github.com/simiancreative/simiango/examples/services/crud"
	_ "github.com/simiancreative/simiango/examples/services/crypt"
	_ "github.com/simiancreative/simiango/examples/services/direct"
	_ "github.com/simiancreative/simiango/examples/services/error"
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
	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/messaging/amqp"
	"github.com/simiancreative/simiango/messaging/kafka"
	"github.com/simiancreative/simiango/meta"
	"github.com/simiancreative/simiango/monitoring/sentry"
	"github.com/simiancreative/simiango/server"
	"github.com/simiancreative/simiango/stats/prometheus"
)

func main() {
	config.Enable()
	logger.Enable()

	sentry.Enable()

	logger.Printf("ENV STARTING AS: %v", os.Getenv("APP_ENV"))
	done, exit := meta.CatchSig()

	mssql.Connect()
	mysql.Connect()
	pg.Connect()

	server.New()

	prometheus.Handle()

	server.EnableHealthCheck()
	server.SetCORS()
	server.AddPprof()
	server.AddSentryMiddleware(lib.SentryScope)

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
