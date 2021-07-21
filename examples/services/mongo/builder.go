package mongoexample

import (
	"github.com/simiancreative/simiango/meta"
	"github.com/simiancreative/simiango/service"
)

// Build godoc
// @Summary Submit a configdata request
// @Description Retrieves configdata records
// @Tags Settings
// @ID configdata
// @Accept  json
// @Produce  json
// @Param configdata body configdata true "ConfigData"
// @Success 204 "No Content"
// @Failure 503 {object} errors.ErrorResult "Connectivity Issue"
// @Router /configdata [get]
func Build(
	requestID meta.RequestId,
	rawHeaders service.RawHeaders,
	rawBody service.RawBody,
	rawParams service.RawParams,
) (
	service.TPL, error,
) {
	// create service instance
	service := Service{
		ID: requestID,
	}

	return service, nil
}
