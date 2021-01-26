package configdata

import (
	"github.com/simiancreative/simiango/service"
)

var Config = service.Config{
	Method: "GET",
	Path:   "/configdata",
	Build:  Build,
}
