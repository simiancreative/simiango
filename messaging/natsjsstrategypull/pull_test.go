package natsjsstrategypull_test

import (
	"context"
	"testing"

	"github.com/simiancreative/simiango/messaging/natsjscon"
	"github.com/simiancreative/simiango/messaging/natsjsstrategypull"
	conmock "github.com/simiancreative/simiango/mocks/messaging/natsjscon"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	config := natsjsstrategypull.Config{
		ConsumerName: "test-consumer",
	}
	strategy, err := natsjsstrategypull.New(config)

	assert.NoError(t, err)
	assert.NotNil(t, strategy)
}

func TestSetup(t *testing.T) {
	d := conmock.NewDependencies()

	config := natsjsstrategypull.Config{
		ConsumerName: "test-consumer",
	}
	strategy, err := natsjsstrategypull.New(config)
	assert.NoError(t, err)

	ctx := context.TODO()

	key := natsjscon.CtxKey("stream-name")
	ctx = context.WithValue(ctx, key, "test")
	key = natsjscon.CtxKey("subject")
	ctx = context.WithValue(ctx, key, "test.subject")
	key = natsjscon.CtxKey("logger")
	ctx = context.WithValue(ctx, key, d.Logger)
	key = natsjscon.CtxKey("connection-manager")
	ctx = context.WithValue(ctx, key, d.Connector)

	err = strategy.Setup(ctx)
	assert.NoError(t, err)
}
