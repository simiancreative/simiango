package mssql

import (
	_ "github.com/denisenkom/go-mssqldb"

	"github.com/simiancreative/simiango/data/sql"
)

var Cx sql.ConnX

func init() {
	Cx = sql.Connect("mssql", "SQLSERVER_URL", "SQLSERVER_REQUIRE_CONNECTION")
}
