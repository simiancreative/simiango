package meta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRescue(t *testing.T) {
	defer RescuePanic(Id(), t)

	testpanic()

	assert.NoError(t, nil)
}

func testpanic() {
	panic("hi")
}
