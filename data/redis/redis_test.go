package redis

import (
	"context"
	"testing"
	"time"

	r "github.com/go-redis/redis/v8"

	"github.com/stretchr/testify/assert"
)

var factories = map[string]interface{}{
	"sss": 23,
}

type testclient struct{}

func (c *testclient) Exists(ctx context.Context, key ...string) *r.IntCmd {
	cmd := &r.IntCmd{}
	return cmd
}

func (c *testclient) Set(ctx context.Context, key string, value interface{}, exp time.Duration) *r.StatusCmd {
	cmd := &r.StatusCmd{}
	return cmd
}

func (c *testclient) Get(ctx context.Context, name string) *r.StringCmd {
	cmd := &r.StringCmd{}

	return cmd
}

func TestClientGet(t *testing.T) {
	c := client{
		Ctx: context.Background(),
		C:   &testclient{},
	}

	val, err := c.Get("42")

	assert.Equal(t, "", *val)
	assert.NoError(t, err)
}

func TestClientSet(t *testing.T) {
	c := client{
		Ctx: context.Background(),
		C:   &testclient{},
	}

	err := c.Set("42", "42", 0)

	assert.NoError(t, err)
}

func TestClientExists(t *testing.T) {
	c := client{
		Ctx: context.Background(),
		C:   &testclient{},
	}

	_, err := c.Exists("42")

	assert.NoError(t, err)
}
