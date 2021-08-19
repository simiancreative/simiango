package pg

import (
	_ "github.com/jackc/pgx/stdlib"

	"github.com/simiancreative/simiango/data/sql"
)

var Cx sql.ConnX

func init() {
	Cx = sql.Connect("pgx", "PG_URL", "PG_REQUIRE_CONNECTION")
}
