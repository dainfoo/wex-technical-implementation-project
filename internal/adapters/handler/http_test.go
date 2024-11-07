package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dainfoo/wex-technical-implementation-project/internal/adapters/handler"
	"github.com/dainfoo/wex-technical-implementation-project/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This file contains tests for the HTTP handler functions.
// It uses Testify for assertions, and runs the tests in parallel.

// TestValidateAndCreateTransaction tests the ValidateAndCreateTransaction function. It tests the following scenarios:
//
// 1. Valid Transaction Data.
// 2. Invalid Timestamp.
// 3. Negative AmountInUSD.
func TestValidateAndCreateTransaction(t *testing.T) {
	// Expected values
	transactionValidTransactionData, err := domain.NewTransaction("Valid Description", time.Now().UTC(), 100.0)
	// Stops the test if the expected results are not as expected (probably the business logic changed)
	require.Empty(t, err)

	tests := []struct {
		name           string
		inputData      handler.TransactionDTO
		expectedErrors []error
		expectedResult *domain.Transaction
	}{
		{
			name: "Valid Transaction Data",
			inputData: handler.TransactionDTO{
				Description: "Valid Description",
				Timestamp:   time.Now().UTC().Format(time.RFC3339),
				AmountInUSD: 100.0,
			},
			expectedErrors: []error{},
			expectedResult: transactionValidTransactionData,
		},
		{
			name: "Invalid Timestamp",
			inputData: handler.TransactionDTO{
				Description: "Test Description",
				Timestamp:   "invalid-timestamp",
				AmountInUSD: 100.0,
			},
			expectedErrors: []error{handler.ErrInvalidTimestampFormat},
			expectedResult: nil,
		},
		{
			name: "Negative AmountInUSD",
			inputData: handler.TransactionDTO{
				Description: "Test Description",
				Timestamp:   time.Now().UTC().Format(time.RFC3339),
				AmountInUSD: -10.0,
			},
			expectedErrors: []error{domain.ErrInvalidAmountInUSD},
			expectedResult: nil,
		},
	}

	transactionHandler := handler.TransactionHandler{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, errs := transactionHandler.ValidateAndCreateTransaction(tt.inputData)

			if len(tt.expectedErrors) > 0 {
				require.Len(t, errs, len(tt.expectedErrors))
				for i, expectedError := range tt.expectedErrors {
					assert.ErrorIs(t, errs[i], expectedError)
				}
			} else {
				assert.Empty(t, errs)
			}

			if tt.expectedResult != nil {
				assert.Equal(t, tt.expectedResult.Description, result.Description)
				assert.Equal(t, tt.expectedResult.AmountInUSD, result.AmountInUSD)
			}
		})
	}
}

// TestWriteSuccessResponse tests the WriteSuccessResponse function. It tests the following scenarios:
//
// 1. Success Response with Data.
func TestWriteSuccessResponse(t *testing.T) {
	tests := []struct {
		name           string
		data           interface{}
		statusCode     int
		expectedOutput string
	}{
		{
			name:           "Success Response With Data",
			data:           map[string]string{"id": "12345"},
			statusCode:     http.StatusOK,
			expectedOutput: "{\"data\":{\"id\":\"12345\"}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rr := httptest.NewRecorder()
			handler.WriteSuccessResponse(rr, tt.data, tt.statusCode)
			assert.Equal(t, tt.statusCode, rr.Code)
			assert.JSONEq(t, tt.expectedOutput, rr.Body.String())
		})
	}
}

// TestWriteErrorResponse tests the WriteErrorResponse function. It tests the following scenarios:
//
// 1. Error Response.
func TestWriteErrorResponse(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		message        string
		expectedOutput string
	}{
		{
			name:           "Error Response",
			statusCode:     http.StatusBadRequest,
			message:        "Invalid request",
			expectedOutput: "{\"error\":\"Invalid request\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rr := httptest.NewRecorder()
			handler.WriteErrorResponse(rr, tt.statusCode, tt.message)
			assert.Equal(t, tt.statusCode, rr.Code)
			assert.JSONEq(t, tt.expectedOutput, rr.Body.String())
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
			expectedErrors:  []error{handler.ErrTimestampEmpty},
		},
		{
			name:            "Invalid Timestamp Format (gibberish)",
			timestampString: "not-a-timestamp",
			expectedErrors:  []error{handler.ErrInvalidTimestampFormat},
		},
		{
			name:            "Invalid Timestamp Format (Unix Timestamp)",
			timestampString: "1617181723",
			expectedErrors:  []error{handler.ErrInvalidTimestampFormat},
		},
		{
			name:            "Correct Format But Future Timestamp",
			timestampString: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			expectedErrors:  []error{handler.ErrInvalidTimestamp},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, errs := handler.ParseAndValidateTimestamp(tt.timestampString)

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
			parsedTimestamp, err := handler.ParseISO8601Timestamp(tt.timestampString)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, parsedTimestamp)
			}
		})
	}
}
