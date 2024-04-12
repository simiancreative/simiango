package amqp

import "github.com/simiancreative/simiango/logger"

var session *Session

func Start() {
	for _, config := range services {
		logger.Debug("AMQP: adding service", logger.Fields{
			"key": config.Key,
		})
	}

	NewSession()
}

func Stop() error {
	return session.Close()
}
