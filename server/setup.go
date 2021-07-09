package server

import (
	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/service"
)

func AddService(addedservice service.Config) {
	services = append(services, addedservice)
}

func AddServices(addedservices []service.Config) {
	services = append(services, addedservices...)
}

func InitService(config service.Config) {
	logger.Debug("Server: adding route", logger.Fields{
		"method": config.Method,
		"path":   config.Path,
	})

	router.Handle(config.Method, config.Path, handleService(config))
}
