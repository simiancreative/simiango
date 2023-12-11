package redis_test

import (
	"encoding/json"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/elliotchance/redismock/v8"
	r "github.com/go-redis/redis/v8"
	"github.com/simiancreative/simiango/data/redis"

	"github.com/stretchr/testify/assert"
)

var factories = map[string]interface{}{
	"sss": 23,
}

func init() {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := r.NewClient(&r.Options{
		Addr: mr.Addr(),
	})

	redis.C = redismock.NewNiceMock(client)
}

type marshaller struct{}

func (m *marshaller) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

func TestClientGet(t *testing.T) {
	err := redis.Get("42", &marshaller{})
	assert.Equal(t, err.Error(), "does_not_exist")
}

func TestClientSet(t *testing.T) {
	err := redis.Set("42", "42", 0)

	assert.NoError(t, err)
}

func TestClientExists(t *testing.T) {
	_, err := redis.Exists("42")
	assert.NoError(t, err)
}

func TestClientDel(t *testing.T) {
	err := redis.Set("42", "42", 0)
	err = redis.Del("42")
	assert.NoError(t, err)
}
