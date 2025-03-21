package connection

import (
	"github.com/simiancreative/simiango/logger"
	"github.com/simiancreative/simiango/messaging/natsjscm"
)

var Shared natsjscm.Connector

func init() {
	log := logger.New()

	connectionConfig := natsjscm.ConnectionConfig{
		Logger: log,
		URL:    "nats://localhost:4222",
	}

	connector, err := natsjscm.NewConnectionManager(connectionConfig)
	if err != nil {
		log.Errorf("failed to create connection manager: %v", err)
		panic(err)
	}

	Shared = connector
}
