package assign

import (
	"github.com/simiancreative/simiango/meta"
	"github.com/simiancreative/simiango/service"
)

func Build(
	requestID meta.RequestId,
	rawHeaders service.RawHeaders,
	rawBody service.RawBody,
	rawParams service.RawParams,
) (
	service.MessageTPL, error,
) {

	// create service instance
	s := assignService{
		rawEvent: rawBody,
	}
	service.ParseBody(rawBody, &s.event)
	return &s, nil
}
