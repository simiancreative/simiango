package amqp

import (
	"encoding/json"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	amqp.Publishing

	Exchange string
	Queue    string
	Type     string
	Data     interface{}
}

// Publish will push data onto the queue
func (p Publisher) Publish() error {
	if !session.isReady {
		return errors.New("failed to push: not connected")
	}

	body, _ := json.Marshal(p.Data)

	if err := session.channel.ExchangeDeclare(
		p.Exchange,
		p.Type,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	queue, err := session.channel.QueueDeclare(
		p.Queue,
		true,  // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)

	if err != nil {
		return err
	}

	return session.channel.Publish(
		p.Exchange, // Exchange
		queue.Name, // Routing key
		false,      // Mandatory
		false,      // Immediate
		amqp.Publishing{
			Headers:         p.Headers,
			ContentType:     p.ContentType,
			ContentEncoding: p.ContentEncoding,
			DeliveryMode:    p.DeliveryMode,
			Priority:        p.Priority,
			CorrelationId:   p.CorrelationId,
			ReplyTo:         p.ReplyTo,
			Expiration:      p.Expiration,
			MessageId:       p.MessageId,
			Timestamp:       p.Timestamp,
			Type:            p.Type,
			UserId:          p.UserId,
			AppId:           p.AppId,
			Body:            body,
		},
	)
}
