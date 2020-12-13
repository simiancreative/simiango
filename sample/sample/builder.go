package sample

import (
	"errors"

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
	resource := &sampleResource{}
	service.ParseBody(rawBody, resource)

	// setup context
	decendants, ok := rawParams.Get("decendantsOf")
	if !ok {
		return nil, errors.New("decendants_filter_required")
	}

	id, _ := rawParams.Get("id")
	params := sampleContext{
		ID:         id.Value,
		Decendants: decendants.Values[0],
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
