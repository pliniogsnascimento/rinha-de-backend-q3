package adapter

import (
	"context"
	"fmt"

	pgx "github.com/jackc/pgx/v5"
)

func NewDbConn(options map[string]string) (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/rinha_backend",
		options["user"],
		options["pass"],
		options["host"],
		options["port"],
	))
}
