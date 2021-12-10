package amqp

import (
	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/service"
)

var services service.Collection

func AddService(config service.Config) {
	logger.Debug("Amqp: adding service", logger.Fields{
		"key": config.Key,
	})

	services = append(services, config)
}
