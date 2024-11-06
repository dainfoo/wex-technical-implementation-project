package domain_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/dainfoo/wex-technical-implementation-project/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This file contains tests for the Transaction domain model. It uses Table Driven Tests to test different scenarios.
// It uses Testify for assertions and runs the tests in parallel.

// TestNewTransaction tests the NewTransaction constructor function. It tests the following scenarios:
//
// 1. Valid Transaction.
// 2. Empty Description.
// 3. Description Too Long.
// 4. Negative Amount In USD.
// 5. Correct Format But Future Timestamp.
func TestNewTransaction(t *testing.T) {
	tests := []struct {
		name                string
		description         string
		timestamp           time.Time
		amountInUSD         float64
		expectedErrors      []error
		expectedTransaction *domain.Transaction
	}{
		{
			name:           "Valid Transaction",
			description:    "Valid Transaction",
			timestamp:      time.Now().UTC(),
			amountInUSD:    500.50,
			expectedErrors: []error{},
			expectedTransaction: &domain.Transaction{
				Description: "Valid Transaction",
				Timestamp:   time.Now().UTC(),
				AmountInUSD: new(big.Float).SetPrec(64).SetFloat64(500.50),
			},
		},
		{
			name:           "Empty Description",
			description:    "",
			timestamp:      time.Now().UTC(),
			amountInUSD:    100.0,
			expectedErrors: []error{domain.ErrDescriptionEmpty},
		},
		{
			name:           "Description Too Long",
			description:    "This description is way too long and should trigger a validation error",
			timestamp:      time.Now().UTC(),
			amountInUSD:    250.0,
			expectedErrors: []error{domain.ErrDescriptionTooLong},
		},
		{
			name:           "Negative Amount In USD",
			description:    "Negative Amount In USD",
			timestamp:      time.Now().UTC(),
			amountInUSD:    -50.0,
			expectedErrors: []error{domain.ErrInvalidAmountInUSD},
		},
		{
			name:           "Correct Format But Future Timestamp",
			description:    "Correct Format But Future Timestamp",
			timestamp:      time.Now().Add(24 * time.Hour),
			amountInUSD:    499.0,
			expectedErrors: []error{domain.ErrInvalidTimestamp},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			transaction, errs := domain.NewTransaction(tt.description, tt.timestamp, tt.amountInUSD)

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

			// Check the transaction fields only if no errors where expected
			if len(tt.expectedErrors) == 0 {
				require.NotNil(t, transaction)
				assert.Equal(t, tt.expectedTransaction.Description, transaction.Description)
				assert.True(t, transaction.Timestamp.Before(time.Now().Add(time.Second)))
				assert.Equal(t, tt.expectedTransaction.AmountInUSD.Cmp(transaction.AmountInUSD), 0)
				assert.NotZero(t, transaction.ID)
			} else {
				assert.Nil(t, transaction)
			}
		})
	}
}

// TestValidateDescription tests the ValidateDescription function. It tests the following scenarios:
//
// 1. Valid Description.
// 2. Empty Description.
// 3. Description Too Long.
func TestValidateDescription(t *testing.T) {
	tests := []struct {
		name           string
		description    string
		expectedErrors []error
	}{
		{
			name:           "Valid Description",
			description:    "Valid Description",
			expectedErrors: []error{},
		},
		{
			name:           "Empty Description",
			description:    "",
			expectedErrors: []error{domain.ErrDescriptionEmpty},
		},
		{
			name:           "Description Too Long",
			description:    "This description is way too long and should trigger a validation error",
			expectedErrors: []error{domain.ErrDescriptionTooLong},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			errs := domain.ValidateDescription(tt.description)

			// Check expected errors
			if len(tt.expectedErrors) > 0 {
				require.Len(t, errs, len(tt.expectedErrors))
				for i, expectedError := range tt.expectedErrors {
					assert.ErrorIs(t, errs[i], expectedError)
				}
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}

// TestValidateAmountInUSD tests the ValidateAmountInUSD function. It tests the following scenarios:
//
// 1. Valid Amount In USD.
// 2. Zero Amount In USD.
// 3. Negative Amount In USD.
func TestValidateAmountInUSD(t *testing.T) {
	tests := []struct {
		name           string
		amountInUSD    float64
		expectedErrors []error
	}{
		{
			name:           "Valid Amount In USD",
			amountInUSD:    10.5,
			expectedErrors: []error{},
		},
		{
			name:           "Zero Amount In USD",
			amountInUSD:    0.0,
			expectedErrors: []error{domain.ErrInvalidAmountInUSD},
		},
		{
			name:           "Negative Amount In USD",
			amountInUSD:    -5.0,
			expectedErrors: []error{domain.ErrInvalidAmountInUSD},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			errs := domain.ValidateAmountInUSD(tt.amountInUSD)

			// Check expected errors
			if len(tt.expectedErrors) > 0 {
				require.Len(t, errs, len(tt.expectedErrors))
				for i, expectedError := range tt.expectedErrors {
					assert.ErrorIs(t, errs[i], expectedError)
				}
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}

// TestRoundToTwoDecimalPlaces tests the RoundToTwoDecimalPlaces function. It tests the following scenarios:
//
// 1. Normal Rounding.
// 2. Round Down.
// 3. Exact Two Decimals.
// 4. Negative Rounding.
// 5. Boundary Rounding Up.
func TestRoundToTwoDecimalPlaces(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{
			name:     "Normal Rounding",
			input:    123.456,
			expected: 123.46,
		},
		{
			name:     "Round Down",
			input:    123.454,
			expected: 123.45,
		},
		{
			name:     "Exact Two Decimals",
			input:    123.45,
			expected: 123.45,
		},
		{
			name:     "Negative Rounding",
			input:    -123.456,
			expected: -123.46,
		},
		{
			name:     "Boundary Rounding Up",
			input:    0.005,
			expected: 0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := domain.RoundToTwoDecimalPlaces(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
