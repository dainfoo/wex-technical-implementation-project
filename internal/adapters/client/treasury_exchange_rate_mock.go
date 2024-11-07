package client

import (
	"github.com/dainfoo/wex-technical-implementation-project/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

// MockTreasuryExchangeRateAdapter is a mock implementation of the TreasuryExchangeRateAdapter interface.
type MockTreasuryExchangeRateAdapter struct {
	mock.Mock
}

// GetExchangeRate mocks the GetExchangeRate method of the TreasuryExchangeRateAdapter.
func (m *MockTreasuryExchangeRateAdapter) GetExchangeRate(currency string) (*domain.ExchangeRate, error) {
	args := m.Called(currency)
	// Retrieves the values from the mocked call arguments.
	return args.Get(0).(*domain.ExchangeRate), args.Error(1)
}
