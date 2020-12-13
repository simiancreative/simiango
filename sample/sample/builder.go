package sample

import (
	"github.com/simiancreative/simiango/meta"
	"github.com/simiancreative/simiango/service"
)

func Build(
	requestID meta.RequestId,
	rawBody service.RawBody,
	rawParams service.RawParams,
) (
	service.TPL, error,
) {
	// setup body
	resource := &SampleResource{}
	service.ParseBody(rawBody, resource)

	// setup context
	decendants, _ := rawParams.Get("decendantsOf")
	id, _ := rawParams.Get("id")
	params := SampleContext{
		ID:         id.(string),
		Decendants: decendants.([]string)[0],
	}

	// create service instance
	service := sampleService{
		id:     requestID,
		body:   *resource,
		params: params,
	}

	// return service
	return service, nil
}
