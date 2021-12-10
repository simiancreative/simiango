package amqp

import (
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/simiancreative/simiango/logger"
)

var Connection *amqp.Connection
var Channel *amqp.Channel
var Queue amqp.Queue

func Start() {
	if err := connect(); err != nil {
		logger.Fatal("Amqp Error", logger.Fields{"err": err})
	}

	// create consumer
	if err := createConsumer(); err != nil {
		logger.Fatal("Amqp Error", logger.Fields{"err": err})
	}

	// add services to consumer

	// create publisher

	// keep process open
	select {}
}

func Stop() error {
	tag := os.Getenv("AMQP_CONSUMER_TAG")
	// will close() the deliveries channel
	if err := Channel.Cancel(tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := Connection.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer logger.Info("AMQP shutdown OK", logger.Fields{})

	// wait for handle() to exit
	return <-consumerDone
}
