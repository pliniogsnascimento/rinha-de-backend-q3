package adapter

import (
	"context"

	pgx "github.com/jackc/pgx/v5"
)

func NewDbConn(connStr string) (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), connStr)
}
