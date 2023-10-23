package adapter

import (
	"context"
	"fmt"
	"time"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DbConnOpts struct {
	User                     string
	Password                 string
	Host                     string
	Port                     string
	Database                 string
	DefaultMaxConns          int32
	DefaultMinConns          int32
	DefaultMaxConnLifetime   time.Duration
	DefaultMaxConnIdleTime   time.Duration
	DefaultHealthCheckPeriod time.Duration
	DefaultConnectTimeout    time.Duration
}

func NewDbConn(options map[string]string) (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/rinha_backend",
		options["user"],
		options["pass"],
		options["host"],
		options["port"],
	))
}

func NewDbConnPool(opts *DbConnOpts, logger *zap.SugaredLogger) (*pgxpool.Pool, error) {
	connConfig, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		opts.User,
		opts.Password,
		opts.Host,
		opts.Port,
		opts.Database,
	))

	if err != nil {
		panic(err)
	}

	connConfig.MaxConns = opts.DefaultMaxConns
	connConfig.MinConns = opts.DefaultMinConns
	connConfig.MaxConnLifetime = opts.DefaultMaxConnLifetime
	connConfig.MaxConnIdleTime = opts.DefaultMaxConnIdleTime
	connConfig.HealthCheckPeriod = opts.DefaultHealthCheckPeriod
	connConfig.ConnConfig.ConnectTimeout = opts.DefaultConnectTimeout

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
}
