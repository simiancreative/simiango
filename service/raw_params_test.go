package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRawParams(t *testing.T) {
	params := RawParams{
		ParamItem{
			Key:   "id",
			Value: "12345",
		},
	}

	val, ok := params.Get("id")
	assert.IsType(t, ParamItem{}, val)
	assert.True(t, ok)

	item := struct {
		ID int `param:"id"`
	}{}

	err := params.Assign(&item)

	assert.Equal(t, 12345, item.ID)
	assert.NoError(t, err)
}
