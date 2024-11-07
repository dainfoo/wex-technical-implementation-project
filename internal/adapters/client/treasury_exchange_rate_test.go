package client_test

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/dainfoo/wex-technical-implementation-project/internal/adapters/client"
	"github.com/dainfoo/wex-technical-implementation-project/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// This file contains tests for the Treasury API implementation of the ExchangeRateService interface.
// It uses Testify for assertions and mocking, and runs the tests in parallel.

// TestGetExchangeRates tests the GetExchangeRates method of the TreasuryExchangeRateAdapter.
// It tests the following scenarios:
//
// 1. Successful Response With One Exchange Rate.
// 2. Successful Response With Multiple Exchange Rates.
// 3. Non-200 Response.
// 4. JSON Decoding Error.
// 5. No Data In Response.
func TestGetExchangeRates(t *testing.T) {
	// Expected results
	successfulResponseExchangeRate1, err := domain.NewExchangeRate("Real", 5.434, time.Date(2024, 9, 30, 0, 0, 0, 0, time.UTC))
	// Stops the test if the expected results are not as expected (probably the business logic changed)
	require.Empty(t, err)
	successfulResponseExchangeRate2, err := domain.NewExchangeRate("Real", 5.5, time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC))
	require.Empty(t, err)

	successfulResponseExchangeRates := []*domain.ExchangeRate{successfulResponseExchangeRate1}
	successfulResponseMultipleExchangeRates := []*domain.ExchangeRate{successfulResponseExchangeRate1, successfulResponseExchangeRate2}

	tests := []struct {
		name          string
		mockResponse  *http.Response
		mockError     error
		expectedRates []*domain.ExchangeRate
		expectedError error
	}{
		{
			name: "Successful Response With One Exchange Rate",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"data":[{"currency":"Real","exchange_rate":"5.434","record_calendar_day":"30","record_calendar_month":"09","record_calendar_year":"2024"}]}`)),
			},
			expectedRates: successfulResponseExchangeRates,
			expectedError: nil,
		},
		{
			name: "Successful Response With Multiple Exchange Rates",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"data":[{"currency":"Real","exchange_rate":"5.434","record_calendar_day":"30","record_calendar_month":"09","record_calendar_year":"2024"},{"currency":"Real","exchange_rate":"5.5","record_calendar_day":"30","record_calendar_month":"06","record_calendar_year":"2024"}]}`)),
			},
			expectedRates: successfulResponseMultipleExchangeRates,
			expectedError: nil,
		},
		{
			name: "Non-200 Response",
			mockResponse: &http.Response{
				StatusCode: http.StatusInternalServerError,
			},
			expectedRates: nil,
			expectedError: client.ErrTreasuryAPIResponse,
		},
		{
			name: "JSON Decoding Error",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`invalid json`)),
			},
			expectedRates: nil,
			expectedError: client.ErrDecodingResponse,
		},
		{
			name: "No Data In Response",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"data":[]}`)),
			},
			expectedRates: nil,
			expectedError: client.ErrExchangeRateNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Creates a mock client
			mockClient := new(client.MockTreasuryExchangeRateAdapter)
			mockClient.On("Get", mock.Anything).Return(tt.mockResponse, tt.mockError)

			treasuryAdapter := client.NewConcreteTreasuryExchangeRateAdapter(mockClient)
			actualRates, actualError := treasuryAdapter.GetExchangeRates("Real")

			// Asserts the results
			if tt.expectedError != nil {
				assert.ErrorIs(t, actualError, tt.expectedError)
			} else {
				assert.NoError(t, actualError)
				assert.Equal(t, len(tt.expectedRates), len(actualRates))
				for i, expectedRate := range tt.expectedRates {
					assert.Equal(t, expectedRate.CurrencyName, actualRates[i].CurrencyName)
					assert.Equal(t, expectedRate.Rate.Cmp(actualRates[i].Rate), 0)
					assert.Equal(t, expectedRate.DateOfRecord, actualRates[i].DateOfRecord)
				}
			}

			// Ensure that the mock client was called correctly
			mockClient.AssertExpectations(t)
		})
	}
}
