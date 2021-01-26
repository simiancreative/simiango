package pg

import (
	"context"
	"database/sql"
	"os"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var Ctx context.Context
var C Conn
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

type Conn interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryFunc(context.Context, string, []interface{}, []interface{}, func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	SendBatch(context.Context, *pgx.Batch) pgx.BatchResults
	Close()
}

func init() {
	addr := os.Getenv("PG_URL")

	Ctx = context.Background()
	C, _ = pgxpool.Connect(Ctx, addr)
	Cx, _ = sqlx.Connect("pgx", addr)
}
