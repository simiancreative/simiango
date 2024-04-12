package server

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	servertiming "github.com/p768lwy3/gin-server-timing"

	"github.com/simiancreative/simiango/meta"
	"github.com/simiancreative/simiango/service"
)

func parseRequest(c *gin.Context) service.Req {
	id := meta.Id()
	c.Header("X-Request-ID", string(id))
	c.Set("request_id", string(id))
	timer := servertiming.FromContext(c)

	if hub := sentrygin.GetHubFromContext(c); hub != nil {
		hub.Scope().SetTag("request-id", string(id))
	}

	return service.Req{
		ID:      id,
		Headers: parseHeaders(c.Request),
		Body:    rawBody(c.Request.Body),
		Params:  parseParams(c.Params, c.Request.URL),
		Timer:   timer,
		Context: c,
	}
}
