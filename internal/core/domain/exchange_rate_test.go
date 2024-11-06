package domain_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/dainfoo/wex-technical-implementation-project/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This file contains tests for the exchange rate domain model. It uses Table Driven Tests to test different scenarios.
// It uses Testify for assertions and runs the tests in parallel.

// TestNewExchangeRate tests the NewExchangeRate constructor function. It tests the following scenarios:
//
// 1. Valid Exchange Rate.
// 2. Empty Currency.
// 3. Negative Rate.
// 4. Rate Zero.
// 5. Future Date.
// 6. Valid Exchange Rate With Current Date.
func TestNewExchangeRate(t *testing.T) {
	tests := []struct {
		name                 string
		currency             string
		rate                 float64
		date                 time.Time
		expectedErrors       []error
		expectedExchangeRate *domain.ExchangeRate
	}{
		{
			name:           "Valid Exchange Rate",
			currency:       "Brazil-Real",
			rate:           5.434,
			date:           time.Now(),
			expectedErrors: []error{},
			expectedExchangeRate: &domain.ExchangeRate{
				Currency: "Brazil-Real",
				Rate:     new(big.Float).SetFloat64(5.434),
				Date:     time.Now(),
			},
		},
		{
			name:           "Empty Currency",
			currency:       "",
			rate:           1.2,
			date:           time.Now(),
			expectedErrors: []error{domain.ErrCurrencyEmpty},
		},
		{
			name:           "Negative Rate",
			currency:       "Brazil-Real",
			rate:           -5.434,
			date:           time.Now(),
			expectedErrors: []error{domain.ErrInvalidExchangeRate},
		},
		{
			name:           "Rate Zero",
			currency:       "Brazil-Real",
			rate:           0,
			date:           time.Now(),
			expectedErrors: []error{domain.ErrInvalidExchangeRate},
		},
		{
			name:           "Future Date",
			currency:       "Brazil-Real",
			rate:           5.434,
			date:           time.Now().Add(24 * time.Hour),
			expectedErrors: []error{domain.ErrInvalidDate},
		},
		{
			name:           "Valid Exchange Rate With Current Date",
			currency:       "Brazil-Real",
			rate:           5.434,
			date:           time.Now(),
			expectedErrors: []error{},
			expectedExchangeRate: &domain.ExchangeRate{
				Currency: "Brazil-Real",
				Rate:     new(big.Float).SetFloat64(5.434),
				Date:     time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			exchangeRate, errs := domain.NewExchangeRate(tt.currency, tt.rate, tt.date)

			// Check expected errors
			if len(tt.expectedErrors) > 0 {
				require.NotNil(t, errs)
				require.Len(t, errs, len(tt.expectedErrors))
				for i, expectedError := range tt.expectedErrors {
					assert.ErrorIs(t, errs[i], expectedError)
				}
			} else {
				assert.Empty(t, errs)
			}

			// Check the exchange rate fields only if no errors where expected
			if len(tt.expectedErrors) == 0 {
				require.NotNil(t, exchangeRate)
				assert.Equal(t, tt.expectedExchangeRate.Currency, exchangeRate.Currency)
				assert.Equal(t, tt.expectedExchangeRate.Rate.Cmp(exchangeRate.Rate), 0)
				assert.True(t, exchangeRate.Date.Before(time.Now().Add(time.Second)))
			} else {
				assert.Nil(t, exchangeRate)
			}
		})
	}
}

// TestValidateExchangeRate tests the ValidateExchangeRate function. It tests the following scenarios:
//
// 1. Valid Exchange Rate.
// 2. Empty Currency.
// 3. Negative Rate.
// 4. Rate Zero.
// 5. Future Date.
func TestValidateExchangeRate(t *testing.T) {
	tests := []struct {
		name           string
		currency       string
		rate           float64
		date           time.Time
		expectedErrors []error
	}{
		{
			name:           "Valid Exchange Rate",
			currency:       "United Kingdom-Pound",
			rate:           0.745,
			date:           time.Now(),
			expectedErrors: []error{},
		},
		{
			name:           "Empty Currency",
			currency:       "",
			rate:           1.2,
			date:           time.Now(),
			expectedErrors: []error{domain.ErrCurrencyEmpty},
		},
		{
			name:           "Negative Rate",
			currency:       "United Kingdom-Pound",
			rate:           -0.745,
			date:           time.Now(),
			expectedErrors: []error{domain.ErrInvalidExchangeRate},
		},
		{
			name:           "Rate Zero",
			currency:       "United Kingdom-Pound",
			rate:           0,
			date:           time.Now(),
			expectedErrors: []error{domain.ErrInvalidExchangeRate},
		},
		{
			name:           "Future Date",
			currency:       "United Kingdom-Pound",
			rate:           0.745,
			date:           time.Now().Add(24 * time.Hour),
			expectedErrors: []error{domain.ErrInvalidDate},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			errs := domain.ValidateExchangeRate(tt.currency, tt.rate, tt.date)

			// Check expected errors
			assert.ElementsMatch(t, tt.expectedErrors, errs)
		})
	}
}
