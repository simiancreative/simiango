package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	val := os.Getenv("APP_ENV")
	assert.Equal(t, "dev", val)
}
