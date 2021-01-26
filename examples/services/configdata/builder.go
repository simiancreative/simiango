package configdata

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
	rawBody service.RawBody,
	rawParams service.RawParams,
) (
	service.TPL, error,
) {
	// create service instance
	service := configDataService{
		id:   requestID,
	}

	return service, nil
}
