package server

import (
	"fmt"

	"github.com/simiancreative/simiango/service"
)

func handleDirect(config service.Config, req service.Req) (interface{}, *service.ResultError) {
	if err := handleDirectAuth(req, config); err != nil {
		return nil, handleError(err)
	}

	for _, handler := range config.Before {
		if err := handler(config, req); err != nil {
			return nil, handleError(err)
		}
	}

	if err := parseBody(config, &req); err != nil {
		return nil, handleError(err)
	}

	result, err := config.Direct(req)

	if err == nil && result == nil {
		return nil, nil
	}

	if err == nil {
		return result, nil
	}

	return nil, handleError(err)
}

func handleDirectAuth(
	req service.Req,
	config service.Config,
) error {
	var err error

	if !config.IsPrivate {
		return nil
	}

	err = config.Auth(req)
	if err != nil {
		err = fmt.Errorf("Authentication Failed: %v", err.Error())
		resultErr := service.ToResultError(err, "service auth failed", 401)
		return resultErr
	}

	return nil
}
