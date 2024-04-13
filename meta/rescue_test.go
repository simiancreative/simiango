package meta

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/simiancreative/simiango/logger"
)

func TestRescue(t *testing.T) {
	logger.Enable()

	defer RescuePanic(Id(), t)

	testpanic()

	assert.NoError(t, nil)
}

func testpanic() {
	panic("hi")
}
