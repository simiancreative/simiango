package redis

import (
	"context"
	"encoding"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	r "github.com/go-redis/redis/v8"
)

var Ctx context.Context
var C r.Cmdable

var Exists ExistsFunc = ExistsDefault
var SetNx SetNXFunc = SetNXDefault
var Set SetFunc = SetDefault
var Get GetFunc = GetDefault
var Del DelFunc = DelDefault

type ExistsFunc func(string) (*bool, error)
type SetNXFunc func(string, interface{}, time.Duration) error
type SetFunc func(string, interface{}, time.Duration) error
type GetFunc func(string, encoding.BinaryUnmarshaler) error
type DelFunc func(string) error

func ExistsDefault(key string) (*bool, error) {
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

func SetNXDefault(key string, value interface{}, exp time.Duration) error {
	return C.SetNX(Ctx, key, value, exp).Err()
}

func SetDefault(key string, value interface{}, exp time.Duration) error {
	return C.Set(Ctx, key, value, exp).Err()
}

func DelDefault(key string) error {
	return C.Del(Ctx, key).Err()
}

func GetDefault(name string, rec encoding.BinaryUnmarshaler) error {
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

	sentinels, ok := os.LookupEnv("REDIS_SENTINELS")
	if ok {
		sentinelList := strings.Split(sentinels, ",")

		fmt.Println("REDIS_SENTINELS configured")

		sentinelPassword, ok := os.LookupEnv("REDIS_SENTINELS_PASS")
		if !ok {
			sentinelPassword = ""
		}
		masterName, ok := os.LookupEnv("REDIS_MASTER_NAME")
		if !ok {
			masterName = ""
		}

		c := r.NewFailoverClient(&r.FailoverOptions{
			MasterName:       masterName,
			SentinelAddrs:    sentinelList,
			SentinelPassword: sentinelPassword,
			Password:         pass, // no password set
			DB:               0,    // use default DB
		})
		C = c
	}

	if !ok {
		c := r.NewClient(&r.Options{
			Addr:     addr,
			Password: pass, // no password set
			DB:       0,    // use default DB
		})
		C = c
	}

	Ctx = context.Background()
}
