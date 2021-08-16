package ll

import (
	"github.com/simiancreative/simiango/service"
)

var Config = service.Config{
	IsPrivate: false,
	Method:    "ll",
	Path:      "/ll",
	Build:     Build,
}