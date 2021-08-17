package mysql

// keep these organized by category with an empty line between each
// 1. core
// 2. remote
// 3. local
import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/simiancreative/simiango/logger"

	"github.com/jmoiron/sqlx"
	sqlLogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/logrusadapter"
	"github.com/sirupsen/logrus"
)

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
	addr := os.Getenv("MYSQL_URL")

	dd, _ := sql.Open("mysql", addr)

	if logger.Level() >= logrus.DebugLevel {
		logger := logger.New()
		loggerAdapter := logrusadapter.New(logger)

		dd = sqlLogger.OpenDriver(
			addr,
			dd.Driver(),
			loggerAdapter,
		)
	}

	_, mustConnect := os.LookupEnv("MYSQL_REQUIRE_CONNECTION")
	err := dd.Ping()
	if mustConnect && err != nil {
		panic(err)
	}

	Cx = sqlx.NewDb(dd, "mysql")
}
