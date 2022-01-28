package kafka

import (
	"github.com/simiancreative/simiango/meta"
	"github.com/simiancreative/simiango/service"
)

func buildService(
	requestID meta.RequestId,
	config service.Config,
	body service.RawBody,
) (service.TPL, error) {
	s, err := config.Build(requestID, service.RawHeaders{}, body, service.RawParams{})
	if err == nil {
		return s, nil
	}

	return nil, err
}
