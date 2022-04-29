package kafka

import (
	"github.com/simiancreative/simiango/service"
)

var services service.Collection

func AddService(config service.Config) {
	kl.Printf("adding service, key: %s", config.Key)
	services = append(services, config)
}
