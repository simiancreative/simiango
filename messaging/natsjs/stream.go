package natsjs

import (
	"crypto/tls"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type conn struct {
	nc *nats.Conn
	js jetstream.JetStream
}

var reconnect = true
var connection *conn
var timeout = 30 * time.Second

func SetTimeout(t time.Duration) {
	timeout = t * time.Second
}

func SetReconnect(r bool) {
	reconnect = r
}

// Connect establishes a connection to the NATS server.
//
// USAGE:
//
// import "api/lib/stream"
//
// stream.New().
//
//	InitStream("default").
//	NewMessage().
//	SetSubject("notifications.email", params.Kind, params.Category).
//	SetData(map[string]int{"pending_notifications_id": row.ID}).
//	Publish()
func Connect() {
	if connection != nil && connection.nc != nil && !connection.nc.IsClosed() {
		return
	}

	host := os.Getenv("NATS_HOST")
	if host == "" {
		panic("NATS_HOST is not set")
	}

	opts := nats.Options{
		Url: host,

		AllowReconnect:       reconnect,
		RetryOnFailedConnect: true,
		MaxReconnect:         10,
		ReconnectWait:        10 * time.Second,

		// this is acceptable because the nats server is not exposed outside the cluster
		// it is currently using a self signed certificate so we need to skip the verification
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	nc, err := opts.Connect()
	if err != nil {
		panic(err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		panic(err)
	}

	connection = &conn{nc, js}
}

func Close() {
	if connection == nil || connection.nc == nil {
		return
	}

	connection.nc.Close()
	connection = nil
}

func New() *Client {
	return &Client{}
}
