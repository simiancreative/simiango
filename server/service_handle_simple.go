package server

import (
	"fmt"

	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/service"
)

func handleDirect(config service.Config, req service.Req) (interface{}, *service.ResultError) {
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
		logger.Error("Authentication Failed", logger.Fields{"err": err})
		resultErr := service.ToResultError(err, "service auth failed", 401)
		return resultErr
	}

	return nil
}
