package services

import (
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

// FindTransactionAndExchangeRateFromCurrency returns a transaction and exchange rate from a given currency name.
func (ts *TransactionService) FindTransactionAndExchangeRateFromCurrency(id uuid.UUID, currencyName string) (*domain.Transaction, *domain.ExchangeRate, error) {
	log.Info().Str("transaction_id", id.String()).Str("currency_name", currencyName).Msg("retrieving transaction and exchange rate")
	transaction, err := ts.transactionRepository.FindTransaction(id)
	if err != nil {
		return nil, nil, err
	}
	exchangeRate, err := ts.exchangeRateAdapter.GetExchangeRate(currencyName)
	if err != nil {
		return nil, nil, err
	}
	return transaction, exchangeRate, nil
}
