package server

import (
	"fmt"

	"github.com/simiancreative/simiango/service"
)

func handleDefault(config service.Config, req service.Req) (interface{}, *service.ResultError) {
	s, err := buildService(config, req)
	if err != nil {
		return nil, handleError(err)
	}

	return serviceResult(s)
}

func buildService(
	config service.Config,
	req service.Req,
) (service.TPL, *service.ResultError) {
	requestID := req.ID
	parsedHeaders := req.Headers
	parsedBody := req.Body
	parsedParams := req.Params

	s, err := config.Build(requestID, parsedHeaders, parsedBody, parsedParams)
	if err != nil {
		return nil, handleError(err)
	}

	err = handleAuth(req, s, config)

	if err == nil {
		return s, nil
	}

	for _, handler := range config.Before {
		if err := handler(config, req); err != nil {
			return nil, handleError(err)
		}
	}

	return nil, handleError(err)
}

func handleAuth(
	req service.Req,
	s service.TPL,
	config service.Config,
) error {
	var err error

	if !config.IsPrivate {
		return nil
	}

	if _, ok := interface{}(s).(service.PrivateTPL); config.IsPrivate && !ok {
		return fmt.Errorf("Private Service requires the Auth method")
	}

	ps, _ := interface{}(s).(service.PrivateTPL)
	err = ps.Auth(req.ID, req.Headers, req.Body, req.Params)

	if err != nil {
		err = fmt.Errorf("Authentication Failed: %v", err.Error())
		resultErr := service.ToResultError(err, "service auth failed", 401)
		return resultErr
	}

	return nil
}

func serviceResult(s service.TPL) (interface{}, *service.ResultError) {
	result, err := s.Result()

	if err == nil && result == nil {
		return nil, nil
	}

	if err == nil {
		return result, nil
	}

	return nil, handleError(err)
}
