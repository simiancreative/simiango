package amqp

import (
	"fmt"
	"os"

	"github.com/simiancreative/simiango/logger"
)

var consumerDone chan error

func createConsumer() error {
	consumerDone = make(chan error)
	tag := os.Getenv("AMQP_CONSUMER_TAG")

	logger.Debug(
		"Queue bound to Exchange, starting Consumer",
		logger.Fields{"consumer tag": tag},
	)

	deliveries, err := Channel.Consume(
		Queue.Name, // name
		tag,        // consumerTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue Consume: %s", err)
	}

	go handle(deliveries, consumerDone)

	return nil
}
