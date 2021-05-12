package mssql

// keep these organized by category with an empty line between each
// 1. core
// 2. remote
// 3. local
import (
	"context"
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var Ctx context.Context
var Cx ConnX

type ConnX interface {
	sqlx.Queryer

	MapperFunc(func(string) string)
	Rebind(string) string
	Unsafe() *sqlx.DB
	BindNamed(string, interface{}) (string, []interface{}, error)
	NamedQuery(string, interface{}) (*sqlx.Rows, error)
	NamedExec(string, interface{}) (sql.Result, error)
	Select(interface{}, string, ...interface{}) error
	Get(interface{}, string, ...interface{}) error
	MustBegin() *sqlx.Tx
	Beginx() (*sqlx.Tx, error)
	MustExec(string, ...interface{}) sql.Result
	Preparex(string) (*sqlx.Stmt, error)
	PrepareNamed(string) (*sqlx.NamedStmt, error)
}

func init() {
	_, mustConnect := os.LookupEnv("MYSQL_REQUIRE_CONNECTION")
	addr := os.Getenv("MYSQL_URL")

	Ctx = context.Background()

	if !mustConnect {
		Cx, _ = sqlx.Connect("mysql", addr)
	}

	if mustConnect {
		Cx = sqlx.MustConnect("mysql", addr)
	}
}
