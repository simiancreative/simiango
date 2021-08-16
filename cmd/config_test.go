package ll

import (
	"testing"

	"github.com/simiancreative/simiango/service"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	assert.Equal(t, Config.Path, "/ll")
	assert.Equal(t, Config.Method, "ll")
	assert.IsType(t, Config.Build, service.Config{}.Build)
}