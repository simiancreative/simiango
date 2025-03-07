package natsjscon

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/simiancreative/simiango/logger"
)

// worker processes messages using the configured strategy
func (c *Consumer) worker(id int) {
	defer c.wg.Done()

	c.logger.Debug("starting worker", logger.Fields{
		"worker_id": id,
		"consumer":  c.config.ConsumerName,
	})

	for {
		select {
		case <-c.ctx.Done():
			c.logger.Debug("worker stopping", logger.Fields{
				"worker_id": id,
				"reason":    "context cancelled",
			})

			return
		default:
			err := c.consumeAndProcess(id)

			if err == nil {
				continue
			}

			if errors.Is(err, context.Canceled) {
				continue
			}

			c.logger.Debug("error in worker", logger.Fields{
				"worker_id": id,
				"error":     err.Error(),
			})

			// Brief backoff on error
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// consumeAndProcess uses the strategy to get messages and processes them
func (c *Consumer) consumeAndProcess(workerID int) error {
	// Use strategy to consume messages
	msgs, err := c.strategy.Consume(c.ctx, workerID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			// This is normal when no messages are available
			return nil
		}

		return fmt.Errorf("failed to consume messages: %w", err)
	}

	if len(msgs) == 0 {
		return nil
	}

	c.logger.Debug("processing batch", logger.Fields{
		"worker_id":  workerID,
		"batch_size": len(msgs),
	})

	// Create processing context with timeout
	procCtx, procCancel := context.WithTimeout(c.ctx, c.config.ProcessTimeout)
	defer procCancel()

	// Process batch
	c.mu.RLock()
	processor := c.processor
	c.mu.RUnlock()

	// Call the batch processor
	results := processor(procCtx, msgs)

	// Handle results
	for msg, status := range results {
		c.handleResult(msg, status)
	}

	return nil
}
