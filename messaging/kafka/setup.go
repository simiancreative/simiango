package kafka

import (
	"github.com/simiancreative/simiango/service"
)

var services service.Collection

func AddService(config service.Config) {
	services = append(services, config)
}
