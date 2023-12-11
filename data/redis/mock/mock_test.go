package redismock_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/simiancreative/simiango/data/redis"
)

func TestRegister(t *testing.T) {
	err := redis.Set("42", "42", 0)

	assert.NoError(t, err)
}
