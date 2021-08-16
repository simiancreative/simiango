package ll

import (
	"testing"

	"github.com/simiancreative/simiango/meta"
	"github.com/simiancreative/simiango/service"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	requestID := meta.Id()
	rawBody := service.RawBody("")
	rawParams := service.RawParams{}
	rawHeaders := service.RawHeaders{}

	service, err := Build(requestID, rawHeaders, rawBody, rawParams)

	assert.NoError(t, err)
	assert.IsType(t, service, &llService{})
}
