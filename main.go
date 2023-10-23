package main

import (
	"fmt"
	"os"

	"github.com/pliniogsnascimento/rinha-de-backend-q3/pkg/adapter"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	dbDefaults = map[string]string{
		"user": "postgres",
		"pass": "example",
		"host": "localhost",
		"port": "5432",
	}

	serverDefaults = map[string]string{
		"port": "8080",
	}
)

func main() {
	// Logger
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	// Configs
	viper.SetDefault("db", dbDefaults)
	viper.SetDefault("server", serverDefaults)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// DB Connection
	// dbConn, err := adapter.NewDbConn(viper.GetStringMapString("db"))
	dbConn, err := adapter.NewDbConnPool(viper.GetStringMapString("db"), sugar)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer dbConn.Close()

	// Services
	personRepository := adapter.NewPersonRepo(dbConn, sugar)

	// Server
	server := adapter.NewServer(personRepository, sugar, viper.GetStringMapString("server"))
	server.StartServer()
}
