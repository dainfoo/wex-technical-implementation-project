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

// MockHTTPClient is a mock type for HTTPClient interface.
type MockHTTPClient struct {
	mock.Mock
}

// Get is a mock method for HTTPClient interface.
func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	args := m.Called(url)
	return args.Get(0).(*http.Response), args.Error(1)
}

// TestGetExchangeRate tests the GetExchangeRate method of the TreasuryExchangeRateAdapter.
// It tests the following scenarios:
//
// 1. Successful Response.
// 2. Non-200 Response.
// 3. JSON Decoding Error.
// 4. No Data In Response.
func TestGetExchangeRate(t *testing.T) {
	// Expected results
	successfulResponseExchangeRate, err := domain.NewExchangeRate("Real", 5.434, time.Date(2024, 9, 30, 0, 0, 0, 0, time.UTC))
	// Stops the test if the expected results are not as expected (probably the business logic changed)
	require.Empty(t, err)

	tests := []struct {
		name          string
		mockResponse  *http.Response
		mockError     error
		expectedRate  *domain.ExchangeRate
		expectedError error
	}{
		{
			name: "Successful Response",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"data":[{"currency":"Real","exchange_rate":"5.434","record_calendar_day":"30","record_calendar_month":"09","record_calendar_year":"2024"}]}`)),
			},
			expectedRate:  successfulResponseExchangeRate,
			expectedError: nil,
		},
		{
			name: "Non-200 Response",
			mockResponse: &http.Response{
				StatusCode: http.StatusInternalServerError,
			},
			expectedRate:  nil,
			expectedError: client.ErrTreasuryAPIResponse,
		},
		{
			name: "JSON Decoding Error",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`invalid json`)),
			},
			expectedRate:  nil,
			expectedError: client.ErrDecodingResponse,
		},
		{
			name: "No Data In Response",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"data":[]}`)),
			},
			expectedRate:  nil,
			expectedError: client.ErrExchangeRateNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Creates a mock client
			mockClient := new(MockHTTPClient)
			mockClient.On("Get", mock.Anything).Return(tt.mockResponse, tt.mockError)

			treasuryAdapter := client.NewConcreteTreasuryExchangeRateAdapter(mockClient)
			actualRate, actualError := treasuryAdapter.GetExchangeRate("Real")

			// Asserts the results
			if tt.expectedError != nil {
				assert.ErrorIs(t, actualError, tt.expectedError)
			} else {
				assert.NoError(t, actualError)
				assert.Equal(t, tt.expectedRate.CurrencyName, actualRate.CurrencyName)
				assert.Equal(t, tt.expectedRate.Rate.Cmp(actualRate.Rate), 0)
				assert.Equal(t, tt.expectedRate.DateOfRecord, actualRate.DateOfRecord)
			}

			// Ensure that the mock client was called correctly
			mockClient.AssertExpectations(t)
		})
	}
}
