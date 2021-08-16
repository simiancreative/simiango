package ll

import (
	"fmt"
	"testing"

	"github.com/simiancreative/simiango/meta"
	"github.com/simiancreative/simiango/service"
	"github.com/simiancreative/simiango/token"

	"github.com/stretchr/testify/assert"
)

var testTokenStr string

func init() {
	testTokenStr = token.GenWithSecret(token.Claims{
	}, bytes("key"), 0)
}

func TestAuth(t *testing.T) {
	var requestID meta.RequestId = meta.Id()
	var rawBody service.RawBody = []byte("")
	var rawParams = service.RawParams{}
	var rawHeaders = service.RawHeaders{
		{
			Key:    "Authorization",
			Values: []string{fmt.Sprintf("Bearer %s", testTokenStr)},
		},
	}

	service := &llService{}
	ok := service.Auth(requestID, rawHeaders, rawBody, rawParams)

	assert.Equal(t, ok, true)
}