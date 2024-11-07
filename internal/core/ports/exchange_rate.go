package ports

import (
	"github.com/dainfoo/wex-technical-implementation-project/internal/core/domain"
)

// This file contains the ports provided by the business logic to the external world.

// ExchangeRateService is the interface that the business logic provides for any adapter that wants to implement
// exchange rate retrieval.
type ExchangeRateService interface {
	GetExchangeRate(currencyName string) (*domain.ExchangeRate, error)
}
