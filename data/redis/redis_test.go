package redis

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/elliotchance/redismock/v8"
	r "github.com/go-redis/redis/v8"

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

	C = redismock.NewNiceMock(client)
}

func TestClientGet(t *testing.T) {
	_, err := Get("42")
	assert.Equal(t, err.Error(), "does_not_exist")
}

func TestClientSet(t *testing.T) {
	err := Set("42", "42", 0)

	assert.NoError(t, err)
}

func TestClientExists(t *testing.T) {
	_, err := Exists("42")
	assert.NoError(t, err)
}
