package ll

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResult(t *testing.T) {
	s := &llService{}
	result, err := s.Result()

	assert.Equal(t, err, nil)
	assert.Equal(t, result, "")
}