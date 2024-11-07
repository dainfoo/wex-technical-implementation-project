package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dainfoo/wex-technical-implementation-project/internal/adapters/client"
	"github.com/dainfoo/wex-technical-implementation-project/internal/adapters/handler"
	"github.com/dainfoo/wex-technical-implementation-project/internal/adapters/repository"
	"github.com/dainfoo/wex-technical-implementation-project/internal/core/services"
)

func main() {
	// Initializes resources
	transactionRepository, err := repository.NewTransactionRepositoryBoltDB("wex-db", "transactions")
	if err != nil {
		log.Fatalf("the transaction repository creationg failed: %v", err)
	}
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	treasuryExchangeRateConverter := client.NewConcreteTreasuryExchangeRateAdapter(httpClient)
	transactionService := services.NewTransactionService(transactionRepository, treasuryExchangeRateConverter)
	transactionHandler := handler.NewTransactionHandler(*transactionService)
	transactionHandler.StartServer("3000")
}
