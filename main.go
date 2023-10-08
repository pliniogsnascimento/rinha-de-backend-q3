package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pliniogsnascimento/rinha-de-backend-q3/pkg/adapter"
	"go.uber.org/zap"
)

func main() {
	// DB Connection
	dbConn, err := adapter.NewDbConn("postgres://postgres:example@localhost:5432/rinha_backend")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbConn.Close(context.Background())

	// Logger
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	// Services
	personRepository := adapter.NewPersonRepo(dbConn, sugar)

	// Server
	server := adapter.NewServer(personRepository, sugar)
	server.StartServer()
}
