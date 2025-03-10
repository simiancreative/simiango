package natsjscon

import (
	"context"
	"errors"
	"fmt"

	"github.com/simiancreative/simiango/logger"
)

type CtxKey string

// Start begins consuming messages
func (c *Consumer) Start(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	ctx = c.setupCtx(ctx)

	c.debug("setting up consumer")
	if err := c.setup(); err != nil {
		return fmt.Errorf("failed to setup consumer: %w", err)
	}

	c.debug("validating consumer configuration")
	if err := c.validate(); err != nil {
		return fmt.Errorf("invalid consumer configuration: %w", err)
	}

	c.debug("is consumer already running?")
	if c.running {
		return errors.New("consumer is already running")
	}

	// Set up the strategy
	c.debug("setting up consumption strategy")
	if err := c.strategy.Setup(ctx); err != nil {
		return fmt.Errorf("failed to setup consumption strategy: %w", err)
	}

	c.ctx, c.cancel = context.WithCancel(ctx)
	c.running = true

	// Start worker pool
	c.debug("starting worker pool")
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

func (c *Consumer) setupCtx(ctx context.Context) context.Context {
	key := CtxKey("stream-name")
	ctx = context.WithValue(ctx, key, c.config.StreamName)

	key = CtxKey("subject")
	ctx = context.WithValue(ctx, key, c.config.Subject)

	key = CtxKey("logger")
	ctx = context.WithValue(ctx, key, c.logger)

	key = CtxKey("connection-manager")
	ctx = context.WithValue(ctx, key, c.cm)

	return ctx
}
