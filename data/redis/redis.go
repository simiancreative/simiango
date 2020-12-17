package redis

import (
	"context"
	"errors"
	"os"
	"time"

	r "github.com/go-redis/redis/v8"
)

type RedisClient interface {
	Exists(context.Context, ...string) *r.IntCmd
	Get(context.Context, string) *r.StringCmd
	Set(context.Context, string, interface{}, time.Duration) *r.StatusCmd
}

var Client client

type client struct {
	Ctx context.Context
	C   RedisClient
}

func (c *client) Exists(key string) (*bool, error) {
	cmd := c.C.Exists(c.Ctx, key)
	exists := false
	result, err := cmd.Result()

	if err != nil {
		return nil, err
	}

	if result == 1 {
		exists = true
	}

	return &exists, cmd.Err()
}

func (c *client) Set(key string, value interface{}, exp time.Duration) error {
	return c.C.Set(c.Ctx, key, value, exp).Err()
}

func (c *client) Get(name string) (*string, error) {
	val, err := c.C.Get(c.Ctx, name).Result()

	if err == r.Nil {
		return nil, errors.New("does_not_exist")
	}

	if err != nil {
		return nil, err
	}

	return &val, nil
}

func init() {
	addr, ok := os.LookupEnv("REDIS_ADDR")
	if !ok {
		addr = "localhost:6379"
	}

	pass, ok := os.LookupEnv("REDIS_PASS")
	if !ok {
		pass = ""
	}

	c := r.NewClient(&r.Options{
		Addr:     addr,
		Password: pass, // no password set
		DB:       0,    // use default DB
	})

	Client = client{
		Ctx: context.Background(),
		C:   c,
	}
}
