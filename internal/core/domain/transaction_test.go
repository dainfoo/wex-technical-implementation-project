package domain_test

import (
	"math/big"
	"testing"
	"time"

	domain2 "github.com/dainfoo/wex-technical-implementation-project/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This file contains tests for the transaction domain model. It uses Table Driven Tests to test different scenarios.
// It uses Testify for assertions and runs the tests in parallel.

// TestNewTransaction tests the NewTransaction constructor function. It tests the following scenarios:
//
// 1. Valid Transaction.
// 2. Empty Description.
// 3. Description Too Long.
// 4. Negative Amount.
// 5. Empty Timestamp.
// 6. Invalid Timestamp Format (gibberish).
// 7. Invalid Timestamp Format (Unix Timestamp).
// 8. Correct Format But Future Timestamp.
func TestNewTransaction(t *testing.T) {
	tests := []struct {
		name                string
		description         string
		timestampString     string
		amount              float64
		expectedErrors      []error
		expectedTransaction *domain2.Transaction
	}{
		{
			name:            "Valid Transaction",
			description:     "Valid Transaction",
			timestampString: time.Now().UTC().Format(time.RFC3339),
			amount:          100.50,
			expectedErrors:  []error{},
			expectedTransaction: &domain2.Transaction{
				Description: "Valid Transaction",
				AmountInUSD: new(big.Float).SetPrec(64).SetFloat64(100.50),
			},
		},
		{
			name:            "Empty Description",
			description:     "",
			timestampString: time.Now().UTC().Format(time.RFC3339),
			amount:          100.0,
			expectedErrors:  []error{domain2.ErrDescriptionEmpty},
		},
		{
			name:            "Description Too Long",
			description:     "This description is way too long and should trigger a validation error",
			timestampString: time.Now().UTC().Format(time.RFC3339),
			amount:          100.0,
			expectedErrors:  []error{domain2.ErrDescriptionTooLong},
		},
		{
			name:            "Negative Amount",
			description:     "Negative Amount",
			timestampString: time.Now().UTC().Format(time.RFC3339),
			amount:          -50.0,
			expectedErrors:  []error{domain2.ErrInvalidAmount},
		},
		{
			name:            "Empty Timestamp",
			description:     "Empty Timestamp",
			timestampString: "",
			amount:          100.0,
			expectedErrors:  []error{domain2.ErrTimestampEmpty},
		},
		{
			name:            "Invalid Timestamp Format (gibberish)",
			description:     "Invalid Timestamp Format (gibberish)",
			timestampString: "not-a-timestamp",
			amount:          100.0,
			expectedErrors:  []error{domain2.ErrInvalidTimestampFormat},
		},
		{
			name:            "Invalid Timestamp Format (Unix Timestamp)",
			description:     "Invalid Timestamp Format (Unix Timestamp)",
			timestampString: "1617181723",
			amount:          100.0,
			expectedErrors:  []error{domain2.ErrInvalidTimestampFormat},
		},
		{
			name:            "Correct Format But Future Timestamp",
			description:     "Correct Format But Future Timestamp",
			timestampString: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			amount:          100.0,
			expectedErrors:  []error{domain2.ErrInvalidTimestamp},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			transaction, errs := domain2.NewTransaction(tt.description, tt.timestampString, tt.amount)

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
				assert.Equal(t, tt.expectedTransaction.AmountInUSD.Cmp(transaction.AmountInUSD), 0)
				assert.NotZero(t, transaction.ID)
				assert.True(t, transaction.Timestamp.Before(time.Now().Add(time.Second)))
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
			expectedErrors: []error{domain2.ErrDescriptionEmpty},
		},
		{
			name:           "Description Too Long",
			description:    "This description is way too long and should trigger a validation error",
			expectedErrors: []error{domain2.ErrDescriptionTooLong},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			errs := domain2.ValidateDescription(tt.description)

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

// TestValidateAmount tests the ValidateAmount function. It tests the following scenarios:
//
// 1. Valid Amount.
// 2. Zero Amount.
// 3. Negative Amount.
func TestValidateAmount(t *testing.T) {
	tests := []struct {
		name           string
		amount         float64
		expectedErrors []error
	}{
		{
			name:           "Valid Amount",
			amount:         10.5,
			expectedErrors: []error{},
		},
		{
			name:           "Zero Amount",
			amount:         0.0,
			expectedErrors: []error{domain2.ErrInvalidAmount},
		},
		{
			name:           "Negative Amount",
			amount:         -5.0,
			expectedErrors: []error{domain2.ErrInvalidAmount},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			errs := domain2.ValidateAmount(tt.amount)

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

// TestParseAndValidateTimestamp tests the ParseAndValidateTimestamp function. It tests the following scenarios:
//
// 1. Valid Timestamp.
// 2. Empty Timestamp.
// 3. Invalid Timestamp Format (gibberish).
// 4. Invalid Timestamp Format (Unix Timestamp).
// 5. Correct Format But Future Timestamp.
func TestParseAndValidateTimestamp(t *testing.T) {
	tests := []struct {
		name            string
		timestampString string
		expectedErrors  []error
	}{
		{
			name:            "Valid Timestamp",
			timestampString: time.Now().UTC().Format(time.RFC3339),
			expectedErrors:  []error{},
		},
		{
			name:            "Empty Timestamp",
			timestampString: "",
			expectedErrors:  []error{domain2.ErrTimestampEmpty},
		},
		{
			name:            "Invalid Timestamp Format (gibberish)",
			timestampString: "not-a-timestamp",
			expectedErrors:  []error{domain2.ErrInvalidTimestampFormat},
		},
		{
			name:            "Invalid Timestamp Format (Unix Timestamp)",
			timestampString: "1617181723",
			expectedErrors:  []error{domain2.ErrInvalidTimestampFormat},
		},
		{
			name:            "Correct Format But Future Timestamp",
			timestampString: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			expectedErrors:  []error{domain2.ErrInvalidTimestamp},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, errs := domain2.ParseAndValidateTimestamp(tt.timestampString)

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

// TestParseISO8601Timestamp tests the ParseISO8601Timestamp function. It tests the following scenarios:
//
// 1. Valid Timestamp.
// 2. Empty Timestamp.
// 3. Invalid Timestamp Format (Unix Timestamp).
func TestParseISO8601Timestamp(t *testing.T) {
	tests := []struct {
		name            string
		timestampString string
		expected        time.Time
		expectError     bool
	}{
		{
			name:            "Valid Timestamp",
			timestampString: "2023-01-01T12:00:00Z",
			expected:        time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			expectError:     false,
		},
		{
			name:            "Empty Timestamp",
			timestampString: "",
			expectError:     true,
		},
		{
			name:            "Invalid Timestamp Format (Unix Timestamp)",
			timestampString: "1617181723",
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			parsedTimestamp, err := domain2.ParseISO8601Timestamp(tt.timestampString)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, parsedTimestamp)
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
			actual := domain2.RoundToTwoDecimalPlaces(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
