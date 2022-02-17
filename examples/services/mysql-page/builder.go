package mysqlpage

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
	service.TPL, error,
) {
	// create service instance
	service := Service{}

	if err := rawParams.Assign(&service); err != nil {
		return nil, err
	}

	return service, nil
}
