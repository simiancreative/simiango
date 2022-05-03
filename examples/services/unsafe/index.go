package unsafe

import (
	"github.com/simiancreative/simiango/server"
	"github.com/simiancreative/simiango/service"

	m "github.com/simiancreative/simiango/data/mysql"
	"github.com/simiancreative/simiango/data/sql/unsafe"
)

var Config = service.Config{
	Kind:   service.DIRECT,
	Method: "GET",
	Path:   "/unsafe",
	Direct: direct,
}

func direct(req service.Req) (interface{}, error) {
	defer req.Timer.NewMetric("no_desc").Start().Stop()
	defer req.Timer.NewMetric("query").WithDesc("desc is optional").Start().Stop()

	u := unsafe.Unsafe{Cx: m.Cx}
	return u.UnsafeSelect(`select 
		CAST("2017-08-29 11:00:00" AS DATETIME) date_time, 
		CAST("2017-08-29" AS DATE) day, 
		cast(54.56 as DECIMAL(15,3)) de, 
		id, 
		name 
	from products`)
}

// dont forget to import your package in your main.go for initialization
// _ "path/to/project/direct"
func init() {
	server.AddService(Config)
}
