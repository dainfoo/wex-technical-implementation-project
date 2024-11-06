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
// 2. Empty Currency Name.
// 3. Negative Rate.
// 4. Rate Zero.
// 5. Future Date Of Record.
// 6. Valid Exchange Rate With Current Date Of Record.
func TestNewExchangeRate(t *testing.T) {
	tests := []struct {
		name                 string
		currencyName         string
		rate                 float64
		dateOfRecord         time.Time
		expectedErrors       []error
		expectedExchangeRate *domain.ExchangeRate
	}{
		{
			name:           "Valid Exchange Rate",
			currencyName:   "Brazil-Real",
			rate:           5.434,
			dateOfRecord:   time.Now(),
			expectedErrors: []error{},
			expectedExchangeRate: &domain.ExchangeRate{
				CurrencyName: "Brazil-Real",
				Rate:         new(big.Float).SetFloat64(5.434),
				DateOfRecord: time.Now(),
			},
		},
		{
			name:           "Empty Currency Name",
			currencyName:   "",
			rate:           1.2,
			dateOfRecord:   time.Now(),
			expectedErrors: []error{domain.ErrCurrencyNameEmpty},
		},
		{
			name:           "Negative Rate",
			currencyName:   "Brazil-Real",
			rate:           -5.434,
			dateOfRecord:   time.Now(),
			expectedErrors: []error{domain.ErrInvalidExchangeRate},
		},
		{
			name:           "Rate Zero",
			currencyName:   "Brazil-Real",
			rate:           0,
			dateOfRecord:   time.Now(),
			expectedErrors: []error{domain.ErrInvalidExchangeRate},
		},
		{
			name:           "Future Date Of Record",
			currencyName:   "Brazil-Real",
			rate:           5.434,
			dateOfRecord:   time.Now().Add(24 * time.Hour),
			expectedErrors: []error{domain.ErrInvalidDateOfRecord},
		},
		{
			name:           "Valid Exchange Rate With Current Date Of Record",
			currencyName:   "Brazil-Real",
			rate:           5.434,
			dateOfRecord:   time.Now(),
			expectedErrors: []error{},
			expectedExchangeRate: &domain.ExchangeRate{
				CurrencyName: "Brazil-Real",
				Rate:         new(big.Float).SetFloat64(5.434),
				DateOfRecord: time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			exchangeRate, errs := domain.NewExchangeRate(tt.currencyName, tt.rate, tt.dateOfRecord)

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
				assert.Equal(t, tt.expectedExchangeRate.CurrencyName, exchangeRate.CurrencyName)
				assert.Equal(t, tt.expectedExchangeRate.Rate.Cmp(exchangeRate.Rate), 0)
				assert.True(t, exchangeRate.DateOfRecord.Before(time.Now().Add(time.Second)))
			} else {
				assert.Nil(t, exchangeRate)
			}
		})
	}
}

// TestValidateExchangeRate tests the ValidateExchangeRate function. It tests the following scenarios:
//
// 1. Valid Exchange Rate.
// 2. Empty Currency Name.
// 3. Negative Rate.
// 4. Rate Zero.
// 5. Future Date Of Record.
func TestValidateExchangeRate(t *testing.T) {
	tests := []struct {
		name           string
		currencyName   string
		rate           float64
		dateOfRecord   time.Time
		expectedErrors []error
	}{
		{
			name:           "Valid Exchange Rate",
			currencyName:   "United Kingdom-Pound",
			rate:           0.745,
			dateOfRecord:   time.Now(),
			expectedErrors: []error{},
		},
		{
			name:           "Empty Currency Name",
			currencyName:   "",
			rate:           1.2,
			dateOfRecord:   time.Now(),
			expectedErrors: []error{domain.ErrCurrencyNameEmpty},
		},
		{
			name:           "Negative Rate",
			currencyName:   "United Kingdom-Pound",
			rate:           -0.745,
			dateOfRecord:   time.Now(),
			expectedErrors: []error{domain.ErrInvalidExchangeRate},
		},
		{
			name:           "Rate Zero",
			currencyName:   "United Kingdom-Pound",
			rate:           0,
			dateOfRecord:   time.Now(),
			expectedErrors: []error{domain.ErrInvalidExchangeRate},
		},
		{
			name:           "Future Date Of Record",
			currencyName:   "United Kingdom-Pound",
			rate:           0.745,
			dateOfRecord:   time.Now().Add(24 * time.Hour),
			expectedErrors: []error{domain.ErrInvalidDateOfRecord},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			errs := domain.ValidateExchangeRate(tt.currencyName, tt.rate, tt.dateOfRecord)

			// Check expected errors
			assert.ElementsMatch(t, tt.expectedErrors, errs)
		})
	}
}
