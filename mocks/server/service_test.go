package server_test

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/simiancreative/simiango/mocks/server"
	"github.com/stretchr/testify/assert"
)

func TestMockServiceConfig(t *testing.T) {
	s := server.
		NewMockServiceConfig().
		SetConfigs(server.ConfigItem{
			Status: 404,
			Value:  "/some-url",
			Resp:   server.MockResp{"error": "not found"},
		}).
		SetURLKey("SOME_KEY").
		SetURLSuffix("/v1").
		StartServer()

	assert.NotNil(t, s)

	resp, _ := resty.New().R().Get(s.URL + "/some-url")

	assert.Equal(t, 404, resp.StatusCode())

	resp, _ = resty.New().R().Get(s.URL + "/another-url")

	assert.Equal(t, 204, resp.StatusCode())

	defer s.Close()
}
