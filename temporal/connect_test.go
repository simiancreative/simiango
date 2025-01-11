package temporal_test

import (
	"os"
	"testing"

	mock "github.com/simiancreative/simiango/mocks/temporal"
	"github.com/simiancreative/simiango/temporal"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	s := mock.MockService()
	defer s.Stop()

	mock.SetHost(s)

	assert.NotPanics(t, func() {
		temporal.Connect()
	}, "Connect should not panic")
}

func TestConnect_NoHost(t *testing.T) {
	os.Unsetenv("TEMPORAL_HOST")

	assert.Panics(t, func() {
		temporal.Connect()
	}, "Connect should panic")
}

func TestConnect_DialError(t *testing.T) {
	os.Setenv("TEMPORAL_HOST", "localhost:1234")

	assert.Panics(t, func() {
		temporal.Connect()
	}, "Connect should panic")
}
