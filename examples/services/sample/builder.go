package sample

import (
	"errors"

	"github.com/simiancreative/simiango/meta"
	"github.com/simiancreative/simiango/service"
	v "github.com/simiancreative/simiango/validate"
)

// Build godoc
// @Summary sample simian go api
// @Description one saple endpoint
// @Tags sample
// @Accept  json
// @Produce  json
// @Param resource body sampleResource true "Sample Resource"
// @Success 200 {object} sampleResp
// @Router /sample/{id} [post]
func Build(
	requestID meta.RequestId,
	rawHeaders service.RawHeaders,
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

	headers := sampleHeaders{}
	for _, header := range rawHeaders {
		headers[header.Key] = header.Value()
	}

	// create service instance
	service := sampleService{
		id:      requestID,
		body:    *resource,
		params:  params,
		headers: headers,
	}

	// return service
	return service, nil
}
