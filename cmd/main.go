package main

import (
	"net/http"
	"os"
	"time"

	"github.com/dainfoo/wex-technical-implementation-project/internal/adapters/client"
	"github.com/dainfoo/wex-technical-implementation-project/internal/adapters/handler"
	"github.com/dainfoo/wex-technical-implementation-project/internal/adapters/repository"
	"github.com/dainfoo/wex-technical-implementation-project/internal/core/services"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	// Check and load .env if needed
	loadEnvIfNeeded()

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		// Default port if none is provided
		serverPort = "3000"
	}

	// Initializes resources
	transactionRepository, err := repository.NewTransactionRepositoryBoltDB("wex-db", "transactions")
	if err != nil {
		log.Fatal().Err(err).Msg("the transaction repository creation failed")
	}
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	treasuryExchangeRateConverter := client.NewConcreteTreasuryExchangeRateAdapter(httpClient)
	transactionService := services.NewTransactionService(transactionRepository, treasuryExchangeRateConverter)
	transactionHandler := handler.NewTransactionHandler(*transactionService)
	transactionHandler.StartServer(serverPort)
}

// loadEnvIfNeeded checks if the SERVER_PORT variable is set and loads the .env file if not.
func loadEnvIfNeeded() {
	// Check if the SERVER_PORT environment variable is set
	if os.Getenv("SERVER_PORT") == "" {
		log.Info().Msg("SERVER_PORT not set, loading .env file...")
		// Load .env file only if SERVER_PORT is not set
		if err := godotenv.Load(); err != nil {
			log.Fatal().Err(err).Msg("error loading .env file")
		} else {
			log.Info().Msg(".env file loaded successfully")
		}
	} else {
		log.Debug().Msg("SERVER_PORT environment variable is already set. Skipping .env loading.")
	}
}
