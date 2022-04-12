package sql

// keep these organized by category with an empty line between each
// 1. core
// 2. remote
// 3. local
import (
	"database/sql"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/simiancreative/simiango/logger"
	sqlLogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/logrusadapter"
	"github.com/sirupsen/logrus"
)

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

func Connect(driver string, addrVar string, mustConnectVar string) *sqlx.DB {
	addr := os.Getenv(addrVar)

	dd, _ := sql.Open(driver, addr)

	if maxOpenStr, ok := os.LookupEnv("DB_MAX_OPEN_CONNECTIONS"); ok {
		maxOpen, _ := strconv.Atoi(maxOpenStr)
		dd.SetMaxOpenConns(maxOpen)
	}

	if maxIdleStr, ok := os.LookupEnv("DB_MAX_IDLE_CONNECTIONS"); ok {
		maxIdle, _ := strconv.Atoi(maxIdleStr)
		dd.SetMaxIdleConns(maxIdle)
	}

	if ttlStr, ok := os.LookupEnv("DB_CONNECTIONS_MAX_LIFETIME_MINUTES"); ok {
		ttl, _ := strconv.Atoi(ttlStr)
		dd.SetConnMaxLifetime(time.Duration(ttl) * time.Minute)
	}

	if logger.Level() >= logrus.TraceLevel {
		logger := logger.New()
		loggerAdapter := logrusadapter.New(logger)

		dd = sqlLogger.OpenDriver(
			addr,
			dd.Driver(),
			loggerAdapter,
		)
	}

	_, mustConnect := os.LookupEnv(mustConnectVar)
	err := dd.Ping()
	if mustConnect && err != nil {
		panic(err)
	}

	return sqlx.NewDb(dd, driver)
}
