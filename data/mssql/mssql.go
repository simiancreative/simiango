package mssql

import (
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
	"context"
	"database/sql"
	"os"
    "log"
    "fmt"
)

var db *sql.DB
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
	addr := os.Getenv("DATABASE_URL")

	Ctx = context.Background()
	Cx = sqlx.MustConnect("sqlserver", addr)

	var err error
    db, err = sql.Open("sqlserver", addr)
    if err != nil {
		log.Fatal("error creating connection pool: ", err.Error())
    }
    
    err = db.PingContext(Ctx)
    if err != nil {
		log.Fatal(err.Error())
		fmt.Printf("mssql not connected!\n")
    }
	fmt.Printf("mssql connected!\n")
}
