package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/simiancreative/simiango/service"
)

func TestParseRequest(t *testing.T) {
	// Create a new gin context with a recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set the "Timing-Context" key in the context
	c.Set("Timing-Context", "some value") // replace "some value" with the appropriate value

	// Initialize the request object in the context
	body := strings.NewReader("request body") // replace "request body" with the appropriate body
	c.Request, _ = http.NewRequest(http.MethodGet, "/", body)

	// Call the function with the test context
	req := parseRequest(c)

	// Assert that the request ID was set correctly
	assert.Equal(t, c.Writer.Header().Get("X-Request-ID"), string(req.ID))

	// Assert that the returned object is of the correct type
	assert.IsType(t, service.Req{}, req)

	// Assert that the returned object has the correct properties
	id, _ := c.Get("request_id")
	assert.Equal(t, string(req.ID), id)
	assert.NotNil(t, req.Context)
}
