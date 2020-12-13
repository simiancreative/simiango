package meta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPing(t *testing.T) {
	res := Id()

	assert.IsType(t, RequestId(""), res)
}
