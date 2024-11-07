package ports

import (
	"github.com/dainfoo/wex-technical-implementation-project/internal/core/domain"
	"github.com/google/uuid"
)

// This file contains the ports provided by the business logic to the external world.

// TransactionRepository is the interface that the business logic provides for any adapter that wants to implement
// data persistence to the transaction model.
type TransactionRepository interface {
	SaveTransaction(transaction domain.Transaction) error
	FindTransaction(id uuid.UUID) (*domain.Transaction, error)
}

// TransactionService is the interface that the business logic provides for any adapter that wants to implement
// user facing transaction saving and retrieval with currency conversion data.
type TransactionService interface {
	SaveTransaction(transaction domain.Transaction) error
	FindTransactionAndExchangeRateFromCurrency(id uuid.UUID, currencyName string) (*domain.Transaction, *domain.ExchangeRate, error)
}
