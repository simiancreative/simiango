package amqp

import (
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func connect() error {
	if err := createConnection(); err != nil {
		return fmt.Errorf("Amqp: %s", err)
	}

	if err := createChannel(); err != nil {
		return fmt.Errorf("Amqp: %s", err)
	}

	if err := createQueue(); err != nil {
		return fmt.Errorf("Amqp: %s", err)
	}

	return nil
}

func createConnection() error {
	amqpURI := os.Getenv("AMQP_URI")

	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return fmt.Errorf("Amqp Dial: %s", err)
	}

	Connection = conn

	return nil
}

func createChannel() error {
	exchange := os.Getenv("AMQP_EXHANGE_NAME")
	exchangeType := os.Getenv("AMQP_EXHANGE_TYPE")

	channel, err := Connection.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	Channel = channel

	if err = channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	return nil
}

func createQueue() error {
	exchange := os.Getenv("AMQP_EXHANGE_NAME")
	queueName := os.Getenv("AMQP_QUEUE_NAME")
	queueKey := os.Getenv("AMQP_QUEUE_KEY")

	queue, err := Channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue Declare: %s", err)
	}

	if err = Channel.QueueBind(
		queue.Name, // name of the queue
		queueKey,   // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}

	Queue = queue

	return nil
}
