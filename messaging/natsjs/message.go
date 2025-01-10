package natsjs

import (
	"context"
	"encoding/json"
	"strings"
)

type Message struct {
	Subject string
	Data    []byte
}

// SetSubject sets the subject of the message. The subject is used to route the
// message to the correct stream.
//
// even though stream name is not required to relate a message to a stream, it
// forces our naming convention and makes it easier to find the stream that the
// message belongs to
func (m *Message) SetSubject(streamName, prefix string, suffix ...string) *Message {
	suffix = append([]string{streamName, prefix}, suffix...)
	m.Subject = strings.ToLower(strings.Join(suffix, "."))

	return m
}

func (m *Message) SetData(raw interface{}) *Message {
	data, err := json.Marshal(raw)
	if err != nil {
		panic(err)
	}

	m.Data = data

	return m
}

func (m *Message) Publish() error {
	if connection == nil {
		panic("Connection is not established")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := connection.js.Publish(ctx, m.Subject, m.Data)
	if err != nil {
		return err
	}

	return nil
}
