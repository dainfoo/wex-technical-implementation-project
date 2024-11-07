package services

import (
	"fmt"

	"github.com/dainfoo/wex-technical-implementation-project/internal/adapters/client"
	"github.com/dainfoo/wex-technical-implementation-project/internal/core/domain"
	"github.com/dainfoo/wex-technical-implementation-project/internal/core/ports"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// This file implements the TransactionService interface and handles the access of external services to the transaction
// repository and exchange rate adapter through a controlled way.

// TransactionService holds the transaction repository and exchange rate adapter.
type TransactionService struct {
	transactionRepository ports.TransactionRepository
	exchangeRateAdapter   client.TreasuryExchangeRateAdapter
}

// NewTransactionService creates a new TransactionService instance.
func NewTransactionService(transactionRepository ports.TransactionRepository, exchangeRateAdapter client.TreasuryExchangeRateAdapter) *TransactionService {
	return &TransactionService{
		transactionRepository: transactionRepository,
		exchangeRateAdapter:   exchangeRateAdapter,
	}
}

// SaveTransaction saves a transaction.
func (ts *TransactionService) SaveTransaction(transaction domain.Transaction) error {
	return ts.transactionRepository.SaveTransaction(transaction)
}

// FindTransactionAndExchangeRateFromCurrency retrieves a transaction along with the exchange rate applicable on the
// purchase date for a given currency name. The exchange rate is considered only if it is found within the past 6
// months from the purchase date.
func (ts *TransactionService) FindTransactionAndExchangeRateFromCurrency(id uuid.UUID, currencyName string) (*domain.Transaction, *domain.ExchangeRate, error) {
	log.Info().Str("transaction_id", id.String()).Str("currency_name", currencyName).Msg("retrieving transaction and exchange rates")

	transaction, err := ts.transactionRepository.FindTransaction(id)
	if err != nil {
		return nil, nil, err
	}
	exchangeRates, err := ts.exchangeRateAdapter.GetExchangeRates(currencyName)
	if err != nil {
		return nil, nil, err
	}

	// Finds the exchange rate closest to the transaction date (within the last 6 months)
	var closestExchangeRate *domain.ExchangeRate
	for _, exchangeRate := range exchangeRates {
		// Ensures the exchange rate is within the last 6 months
		if exchangeRate.DateOfRecord.Before(transaction.Timestamp.AddDate(0, -6, 0)) {
			continue
		}

		// Sets the closest exchange rate to the first one or if it is closer to the transaction date
		if closestExchangeRate == nil || exchangeRate.DateOfRecord.After(transaction.Timestamp) && exchangeRate.DateOfRecord.Before(closestExchangeRate.DateOfRecord) {
			closestExchangeRate = exchangeRate
		}
	}

	// Returns an error ff no exchange rate is found within the last 6 months
	if closestExchangeRate == nil {
		return nil, nil, fmt.Errorf("no exchange rate found within the last 6 months for currency %s", currencyName)
	}

	return transaction, closestExchangeRate, nil
}
