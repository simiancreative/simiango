package rabbit

import (
	"fmt"

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

	fmt.Println(rawHeaders, rawParams)

	// create service instance
	service := rabbitService{}
	return &service, nil
}
