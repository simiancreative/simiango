package nats

import (
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestJetStreamMsg(t *testing.T) {
	subject := "test.subject"
	data := []byte("test data")
	msg := NewJetStreamMsg(subject, data)

	t.Run("Metadata", func(t *testing.T) {
		metadata, err := msg.Metadata()
		assert.Nil(t, metadata)
		assert.Nil(t, err)
	})

	t.Run("Data", func(t *testing.T) {
		assert.Equal(t, data, msg.Data())
	})

	t.Run("Headers", func(t *testing.T) {
		assert.IsType(t, nats.Header{}, msg.Headers())
	})

	t.Run("Subject", func(t *testing.T) {
		assert.Equal(t, subject, msg.Subject())
	})

	t.Run("Reply", func(t *testing.T) {
		assert.Equal(t, "", msg.Reply())
	})

	t.Run("Ack", func(t *testing.T) {
		assert.Nil(t, msg.Ack())
	})

	t.Run("DoubleAck", func(t *testing.T) {
		assert.Nil(t, msg.DoubleAck(nil))
	})

	t.Run("Nak", func(t *testing.T) {
		assert.Nil(t, msg.Nak())
	})

	t.Run("NakWithDelay", func(t *testing.T) {
		assert.Nil(t, msg.NakWithDelay(0))
	})

	t.Run("InProgress", func(t *testing.T) {
		assert.Nil(t, msg.InProgress())
	})

	t.Run("Term", func(t *testing.T) {
		assert.Nil(t, msg.Term())
	})

	t.Run("TermWithReason", func(t *testing.T) {
		assert.Nil(t, msg.TermWithReason(""))
	})

	t.Run("ackReply", func(t *testing.T) {
		assert.Nil(t, msg.ackReply(nil, nil, false, ackOpts{}))
	})

	t.Run("checkReply", func(t *testing.T) {
		assert.Nil(t, msg.checkReply())
	})
}
