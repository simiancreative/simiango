package pg

import (
	_ "github.com/lib/pq"

	"github.com/simiancreative/simiango/data/sql"
)

var Cx sql.ConnX

func init() {
	Cx = sql.Connect("postgres", "PG_URL", "PG_REQUIRE_CONNECTION")
}
