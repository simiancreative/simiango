package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaders(t *testing.T) {
	headers := RawHeaders{
		ParamItem{
			Key:    "X-Companyid",
			Values: []string{"12345"},
		},
	}

	val, ok := headers.Get("X-Companyid")
	assert.IsType(t, ParamItem{}, val)
	assert.True(t, ok)

	item := struct {
		ID int `header:"X-Companyid"`
	}{}

	err := headers.Assign(&item)

	assert.Equal(t, 12345, item.ID)
	assert.NoError(t, err)
}
