package natsjscon

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/simiancreative/simiango/logger"
)

// handleResult handles the result of message processing
func (c *Consumer) handleResult(msg jetstream.Msg, status ProcessStatus) {
	if msg == nil {
		c.logger.Debug("missing jetstream message reference", nil)
		return
	}

	metadata, err := msg.Metadata()
	if err != nil {
		c.logger.Debug("error getting message metadata", logger.Fields{
			"error": err.Error(),
		})
		return
	}

	switch status {
	case Success:
		handleSuccess(msg, c)
		return
	case Failure:
		handleFailure(msg, metadata, c)
		return
	case TerminalFailure:
		handleTerminalError(msg, c)
	}
}

func handleSuccess(msg jetstream.Msg, c *Consumer) {
	// Acknowledge successful processing
	err := msg.Ack()

	if err == nil {
		return
	}

	c.logger.Error("error acknowledging message", logger.Fields{
		"error": err.Error(),
	})
}

func handleFailure(msg jetstream.Msg, metadata *jetstream.MsgMetadata, c *Consumer) {
	// Check if max retries reached
	limitReached := handleMaxRetries(msg, metadata, c)
	if limitReached {
		return
	}

	// Negative ack for retry
	err := msg.Nak()
	if err == nil {
		return
	}

	// Log error
	c.logger.Error("error negative-acknowledging message", logger.Fields{
		"error": err.Error(),
	})
}

func handleMaxRetries(msg jetstream.Msg, metadata *jetstream.MsgMetadata, c *Consumer) bool {
	if metadata == nil {
		return false
	}

	if metadata.NumDelivered < uint64(c.config.MaxRetries) {
		return false
	}

	err := msg.Ack()
	if err != nil {
		c.logger.Debug("error acknowledging message", logger.Fields{
			"error": err.Error(),
		})
	}

	if !c.config.EnableDLQ || c.dlqHandler == nil {
		return true
	}

	publishToDLQ(msg, "max_retries", c)

	return true
}

func handleTerminalError(msg jetstream.Msg, c *Consumer) {
	// Ack to prevent redelivery
	err := msg.Ack()
	if err != nil {
		c.logger.Debug("error acknowledging message", logger.Fields{
			"error": err.Error(),
		})
	}

	// Send to DLQ if enabled
	if !c.config.EnableDLQ || c.dlqHandler == nil {
		return
	}

	// Create a NATS message for the DLQ
	publishToDLQ(msg, "terminal_error", c)
}

func publishToDLQ(original jetstream.Msg, reason string, c *Consumer) {
	if c.dlqHandler == nil {
		return
	}

	msg := &nats.Msg{
		Subject: original.Subject(),
		Data:    original.Data(),
		Header:  original.Headers(),
	}

	err := c.dlqHandler.PublishMessage(msg, reason)
	if err == nil {
		return
	}

	c.logger.Debug("error sending to DLQ", logger.Fields{
		"error": err.Error(),
	})
}
