package temporal_test

import (
	"testing"

	"github.com/simiancreative/simiango/temporal"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	model := temporal.Register("testModel")
	assert.NotNil(t, model)
	assert.Equal(t, "testModel", model.Name)
}

func TestRegister_ErrorHandling(t *testing.T) {
	temporal.Register("newTestModel")
	assert.Panics(t, func() {
		temporal.Register("newTestModel")
	})
}
