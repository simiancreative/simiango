package natsjs

import (
	"context"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

type Client struct{}

func (c *Client) InitStream(streamName string) *Client {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if connection == nil {
		panic("Connection is not established")
	}

	_, err := connection.js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     streamName,
		Subjects: []string{streamName + ".>"},
	})
	if err != nil {
		panic(err)
	}

	return c
}

func (c *Client) NewMessage() *Message {
	return &Message{}
}
