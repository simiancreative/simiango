package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type _s struct{}

func (s _s) Result() (interface{}, error) {
	return response()
}

func response() (*string, error) {
	return nil, nil
}

func TestServiceResult(t *testing.T) {
	result, err := serviceResult(_s{})

	assert.Nil(t, result)
	assert.Nil(t, err)
}
