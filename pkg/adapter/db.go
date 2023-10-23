package adapter

import (
	"context"
	"fmt"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func NewDbConn(options map[string]string) (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/rinha_backend",
		options["user"],
		options["pass"],
		options["host"],
		options["port"],
	))
}

func NewDbConnPool(options map[string]string, logger *zap.SugaredLogger) (*pgxpool.Pool, error) {
	connConfig, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%s/rinha_backend",
		options["user"],
		options["pass"],
		options["host"],
		options["port"],
	))
	if err != nil {
		panic(err)
	}

	// TODO: Refactor the configs usage
	// const defaultMaxConns = int32(4)
	// const defaultMinConns = int32(0)
	// const defaultMaxConnLifetime = time.Hour
	// const defaultMaxConnIdleTime = time.Minute * 30
	// const defaultHealthCheckPeriod = time.Minute
	// const defaultConnectTimeout = time.Second * 5

	// connConfig.MaxConns = defaultMaxConns
	// connConfig.MinConns = defaultMinConns
	// connConfig.MaxConnLifetime = defaultMaxConnLifetime
	// connConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	// connConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	// connConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	connConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		logger.Debugln("Acquiring the connection from db pool")
		return true
	}

	connConfig.AfterRelease = func(c *pgx.Conn) bool {
		logger.Debugln("Connection released from the pool")
		return true
	}

	connConfig.BeforeClose = func(c *pgx.Conn) {
		logger.Debugln("Closed the connection pool to the database!!")
	}

	return pgxpool.NewWithConfig(context.Background(), connConfig)
	// return pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/rinha_backend",
	// 	options["user"],
	// 	options["pass"],
	// 	options["host"],
	// 	options["port"],
	// ))
}
