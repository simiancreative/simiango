package natsjscon_test

import (
	"context"
	"testing"
	"time"

	conmock "github.com/simiancreative/simiango/mocks/messaging/natsjscon"
	"github.com/stretchr/testify/assert"
)

func TestNewConsumer(t *testing.T) {
	d := conmock.NewDependencies()
	c := d.NewConsumer(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := c.Start(ctx)
	assert.NoError(t, err)

	// wait
	time.Sleep(500 * time.Millisecond)

	// does not panic
	assert.NotPanics(t, func() {
		c.Stop()
	}, "Stop should not panic")
}
