package mysql

import (
	// register dialect for goqu
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	// register driver for mysql
	_ "github.com/go-sql-driver/mysql"

	"github.com/simiancreative/simiango/data/sql"
)

var Cx sql.ConnX

func init() {
	Cx = sql.Connect("mysql", "MYSQL_URL", "MYSQL_REQUIRE_CONNECTION")
}
