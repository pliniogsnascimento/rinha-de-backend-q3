package main

import (
	"fmt"
	"os"
	"time"

	db "github.com/pliniogsnascimento/rinha-de-backend-q3/pkg/adapter/database"
	"github.com/pliniogsnascimento/rinha-de-backend-q3/pkg/adapter/http"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	dbDefaults = &db.DbConnOpts{
		DefaultMaxConns:          int32(4),
		DefaultMinConns:          int32(0),
		DefaultMaxConnLifetime:   time.Hour,
		DefaultMaxConnIdleTime:   time.Minute * 30,
		DefaultHealthCheckPeriod: time.Minute,
		DefaultConnectTimeout:    time.Second * 5,
	}

	serverDefaults = &http.ServerOptions{
		Port: 8080,
		RateLimiting: http.RateLimiterOptions{
			Enable: false,
		},
	}
)

func main() {
	// Logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()
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
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	gormDb, err := db.NewGormDb(dbOpts, sugar)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	// Services
	personRepository := db.NewPersonRepo(gormDb, sugar)

	// Server
	serverOpts := serverDefaults
	err = viper.UnmarshalKey("server", serverOpts)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	server := http.NewServer(personRepository, sugar, serverOpts)
	server.StartServer()
}
