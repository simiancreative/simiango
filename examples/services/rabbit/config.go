package rabbit

import (
	"github.com/simiancreative/simiango/messaging/amqp"
	"github.com/simiancreative/simiango/service"
)

var Config = service.Config{
	Key:       "add_user_fob",
	Build:     Build,
}

// dont forget to import your package in your main.go for initialization
// _ "path/to/project/stream"
func init() {
	amqp.AddService(Config)
}
