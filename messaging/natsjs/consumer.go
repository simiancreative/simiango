package natsjs

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/simiancreative/simiango/errors"
	"github.com/simiancreative/simiango/logger"
)

func NewConsumer(name, stream, subject string) *Consumer {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if connection == nil {
		panic("Connection is not established")
	}

	// retrieve consumer handle from a stream
	cons, err := connection.js.CreateOrUpdateConsumer(
		ctx, stream, jetstream.ConsumerConfig{
			FilterSubject: subject,
			Durable:       name,
			Name:          name,
		},
	)
	if err != nil {
		panic(err)
	}

	return &Consumer{cons: cons}
}

type Consumer struct {
	cons jetstream.Consumer
	// batchSize int
}

// func (c *Consumer) BatchConfiguration(size int) *Consumer {
// 	c.batchSize = size
//
// 	return c
// }

func (c Consumer) Consume(callback func(jetstream.Msg) error) (stop func(), err error) {
	cc, err := c.cons.Consume(func(msg jetstream.Msg) {
		err := callback(msg)

		if err != nil {
			logger.Errorf("Error consuming message: %#v", errors.Unwrap(err))
		}

		logger.Warnf("Acking message: %s", msg.Data())
		msg.Ack()
	})

	stop = func() {
		logger.Debugf("Stopping consumer")
		cc.Stop()
	}

	return stop, err
}

// Some thinking to do before we implement a batch consumer.

// The nats side is relatively easy. We need to decide how we want to handle the
// application logic and how we ack messages.

// Do we want to ack messages in the callback or after the callback has been executed?
//
// 1. If we ack messages after the callback we need to make sure that the callback is
//    idempotent on a per message basis.
//
// 2. If we ack messages inside the callback we need to make sure that the messages
//    can be acked out of order in the case where some messages are acked but not
//    others.

// func (c Consumer) BatchConsume(callback func([]string) error) (stop func(), err error) {
// 	var shouldStop bool
//
// 	go func() {
// 		for {
// 			if shouldStop {
// 				break
// 			}
//
// 			batch, err := c.cons.Fetch(c.batchSize)
// 			if err != nil {
// 				log.Printf("Error fetching messages: %v", err)
// 				continue
// 			}
//
// 			msgs, data := collectMessages(batch)
//
// 			if len(msgs) == 0 {
// 				continue
// 			}
//
// 			err = callback(data)
// 			if err != nil {
// 				logger.Errorf("Error consuming message: %v", err)
// 				continue
// 			}
//
// 			for _, msg := range msgs {
// 				msg.Ack()
// 			}
// 		}
// 	}()
//
// 	stop = func() {
// 		logger.Debugf("Stopping consumer")
// 		shouldStop = true
// 	}
//
// 	return stop, err
// }
//
// func collectMessages(batch jetstream.MessageBatch) (msgs []jetstream.Msg, data []string) {
// 	for msg := range batch.Messages() {
// 		data = append(data, string(msg.Data()))
// 		msgs = append(msgs, msg)
// 	}
//
// 	return
// }
