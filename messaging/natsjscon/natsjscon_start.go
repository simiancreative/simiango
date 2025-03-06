package natsjscon

import (
	"context"
	"errors"
	"fmt"

	"github.com/simiancreative/simiango/logger"
)

// Start begins consuming messages
func (c *Consumer) Start(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.validate(); err != nil {
		return fmt.Errorf("invalid consumer configuration: %w", err)
	}

	if c.running {
		return errors.New("consumer is already running")
	}

	// Set up the strategy
	if err := c.strategy.Setup(ctx); err != nil {
		return fmt.Errorf("failed to setup consumption strategy: %w", err)
	}

	c.ctx, c.cancel = context.WithCancel(ctx)
	c.running = true

	// Start worker pool
	for i := 0; i < c.config.WorkerCount; i++ {
		c.wg.Add(1)
		go c.worker(i)
	}

	c.logger.Debug("consumer started", logger.Fields{
		"stream":       c.config.StreamName,
		"consumer":     c.config.ConsumerName,
		"subject":      c.config.Subject,
		"worker_count": c.config.WorkerCount,
	})

	return nil
}
