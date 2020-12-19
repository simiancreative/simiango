package redis

import (
	"context"
	"encoding"
	"errors"
	"os"
	"time"

	r "github.com/go-redis/redis/v8"
)

var Ctx context.Context
var C r.Cmdable

var Exists ExistsFunc = exists
var SetNx SetNXFunc = setNX
var Set SetFunc = set
var Get GetFunc = get

type ExistsFunc func(string) (*bool, error)

func exists(key string) (*bool, error) {
	cmd := C.Exists(Ctx, key)
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

type SetNXFunc func(string, interface{}, time.Duration) error

func setNX(key string, value interface{}, exp time.Duration) error {
	return C.SetNX(Ctx, key, value, exp).Err()
}

type SetFunc func(string, interface{}, time.Duration) error

func set(key string, value interface{}, exp time.Duration) error {
	return C.Set(Ctx, key, value, exp).Err()
}

type GetFunc func(string, encoding.BinaryUnmarshaler) error

func get(name string, rec encoding.BinaryUnmarshaler) error {
	val, err := C.Get(Ctx, name).Result()

	if err == r.Nil {
		return errors.New("does_not_exist")
	}

	if err != nil {
		return err
	}

	rec.UnmarshalBinary([]byte(val))

	return nil
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

	Ctx = context.Background()
	C = c
}
