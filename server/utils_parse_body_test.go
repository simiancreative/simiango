package server

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/simiancreative/simiango/service"
)

func TestParseBody(t *testing.T) {
	t.Run("when config.Input is nil", func(t *testing.T) {
		config := service.Config{}
		req := &service.Req{}

		err := parseBody(config, req)

		assert.Nil(t, err)
		assert.Nil(t, req.Input)
	})

	t.Run("when config.Input is not nil and body parsing is successful", func(t *testing.T) {
		config := service.Config{
			Input: func() interface{} { return &service.Req{} },
		}
		req := &service.Req{
			Body: []byte(`{"key":"value"}`),
		}

		err := parseBody(config, req)

		assert.Nil(t, err)
		assert.NotNil(t, req.Input)
	})

	t.Run("when config.Input is not nil and body parsing fails", func(t *testing.T) {
		config := service.Config{
			Input: func() interface{} { return &service.Req{} },
		}
		req := &service.Req{
			Body: []byte(`{"key":}`), // invalid JSON
		}

		err := parseBody(config, req)

		assert.NotNil(t, err)
		assert.Nil(t, req.Input)
	})
}
