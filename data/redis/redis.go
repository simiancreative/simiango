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

func Exists(key string) (*bool, error) {
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

func SetNX(key string, value interface{}, exp time.Duration) error {
	return C.SetNX(Ctx, key, value, exp).Err()
}

func Set(key string, value interface{}, exp time.Duration) error {
	return C.Set(Ctx, key, value, exp).Err()
}

func Get(name string, rec encoding.BinaryUnmarshaler) error {
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
