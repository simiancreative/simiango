package pgservice

import (
	"github.com/simiancreative/simiango/meta"
	"github.com/simiancreative/simiango/service"
)

// Build godoc
// @Summary sample simian go pg get
// @Description one pg endpoint
// @Tags pg
// @Accept  json
// @Produce  json
// @Success 200 {object} Products
// @Router /pg-example/products [get]
func Build(
	requestID meta.RequestId,
	rawHeaders service.RawHeaders,
	rawBody service.RawBody,
	rawParams service.RawParams,
) (
	service.TPL, error,
) {
	// create service instance
	service := sampleService{id: requestID}

	// return service
	return service, nil
}
