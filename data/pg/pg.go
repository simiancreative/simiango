package pg

import (
	"context"
	"os"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var Ctx context.Context
var C Conn

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
	addr := os.Getenv("DATABASE_URL")
	Ctx = context.Background()
	C, _ = pgxpool.Connect(Ctx, addr)
}
