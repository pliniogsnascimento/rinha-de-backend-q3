package main

import (
	"fmt"
	"os"
	"time"

	"github.com/pliniogsnascimento/rinha-de-backend-q3/pkg/adapter"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	dbDefaults = &adapter.DbConnOpts{
		DefaultMaxConns:          int32(4),
		DefaultMinConns:          int32(0),
		DefaultMaxConnLifetime:   time.Hour,
		DefaultMaxConnIdleTime:   time.Minute * 30,
		DefaultHealthCheckPeriod: time.Minute,
		DefaultConnectTimeout:    time.Second * 5,
	}

	serverDefaults = &adapter.ServerOptions{
		Port: 8080,
		RateLimiting: adapter.RateLimiterOptions{
			Enable: false,
		},
	}
)

func main() {
	// Logger
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	// Configs
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// DB Connection
	dbOpts := dbDefaults
	err = viper.UnmarshalKey("db", dbOpts)
	sugar.Infoln(dbOpts)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	dbConn, err := adapter.NewDbConnPool(dbOpts, sugar)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbConn.Close()

	// Services
	personRepository := adapter.NewPersonRepo(dbConn, sugar)

	// Server
	serverOpts := serverDefaults
	err = viper.UnmarshalKey("server", serverOpts)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	server := adapter.NewServer(personRepository, sugar, serverOpts)
	server.StartServer()
}
