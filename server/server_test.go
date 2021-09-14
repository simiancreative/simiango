package server

import (
	"testing"

	"github.com/simiancreative/simiango/service"
	"github.com/stretchr/testify/assert"
)

type _s struct {
	result interface{}
}

func (s _s) Result() (interface{}, error) {
	return s.result, nil
}

func TestServiceResult(t *testing.T) {
	result, err := serviceResult(_s{result: nil})

	assert.Equal(t, result, nil)
	assert.Equal(t, err, (*service.ResultError)(nil))
}

func TestServiceResultNilPointer(t *testing.T) {
	var result interface{}
	var str *string

	result = str

	result, err := serviceResult(_s{result: result})

	assert.Equal(t, result, (*string)(nil))
	assert.Equal(t, err, (*service.ResultError)(nil))
}
