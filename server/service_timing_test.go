package server

import (
	"testing"

	"github.com/simiancreative/simiango/service"
	"github.com/stretchr/testify/assert"
)

func TestDirectService(t *testing.T) {
	req := service.Req{}

	config := service.Config{
		Kind: service.DIRECT,
		Direct: func(r service.Req) (interface{}, error) {
			return nil, nil
		},
	}

	result, err := handleDirect(config, req)

	assert.Nil(t, result)
	assert.Nil(t, err)
}
