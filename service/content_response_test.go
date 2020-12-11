package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToContentResponse(t *testing.T) {
	content := ToContentResponse([]interface{}{3, 1, 4}, ContentResponseMeta{
		Total: 5000,
		Page:  1,
		Size:  25,
	})

	assert.Equal(t, 200, content.TotalPages)
	assert.Equal(t, true, content.First)
	assert.Equal(t, false, content.Last)
}
