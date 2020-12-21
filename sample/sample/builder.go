package sample

import (
	"errors"

	"github.com/simiancreative/simiango/meta"
	"github.com/simiancreative/simiango/service"
	v "github.com/simiancreative/simiango/validate"
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

	result := v.Validate(resource)
	if result.HasErrors() {
		return nil, result.ResultError()
	}

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
