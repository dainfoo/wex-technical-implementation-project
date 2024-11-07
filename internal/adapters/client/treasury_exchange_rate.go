package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dainfoo/wex-technical-implementation-project/internal/core/domain"
	"github.com/rs/zerolog/log"
)

// This file contains the implementation of the ExchangeRateService	interface using the Treasury API.

// HTTPClient just wraps te http.Client interface to make it easier to mock in tests.
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// Constants for the Treasury API. Change these if the API changes.
const (
	treasuryAPIEndpoint = "https://api.fiscaldata.treasury.gov/services/api/fiscal_service/v1/accounting/od/rates_of_exchange"
	maxRetries          = 3
	retryDelay          = 1 * time.Second
)

// TreasuryExchangeRateAdapter represents a client to the Treasury API.
type TreasuryExchangeRateAdapter struct {
	client      HTTPClient
	apiEndpoint string
}

// NewTreasuryExchangeRateAdapter creates a new TreasuryExchangeRateAdapter with the given HTTPClient.
func NewTreasuryExchangeRateAdapter(client HTTPClient) *TreasuryExchangeRateAdapter {
	return &TreasuryExchangeRateAdapter{
		client:      client,
		apiEndpoint: treasuryAPIEndpoint,
	}
}

// GetExchangeRate retrieves the most recent exchange rate for a currency with input and response validations.
func (a *TreasuryExchangeRateAdapter) GetExchangeRate(currencyName string) (*domain.ExchangeRate, error) {
	apiURL := buildRequestURL(a, currencyName)

	// Retry mechanism
	var resp *http.Response
	var err error
	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err = a.client.Get(apiURL)
		if err == nil {
			break
		}
		if attempt < maxRetries-1 {
			log.Warn().Err(err).Msg("retrying request due to transient network issue")
			time.Sleep(retryDelay)
		}
	}
	if err != nil {
		return nil, ErrNetworkIssue
	}

	return ProcessResponse(resp, currencyName)
}

// buildRequestURL constructs the URL for the Treasury API request.
func buildRequestURL(a *TreasuryExchangeRateAdapter, currencyName string) string {
	return fmt.Sprintf("%s?&sort=-record_date&format=json&page[number]=1&page[size]=1"+
		"&fields=currency,exchange_rate,record_date,record_calendar_day,record_calendar_month,record_calendar_year"+
		"&filter=currency:eq:%s", a.apiEndpoint, url.QueryEscape(currencyName))
}

// ProcessResponse reads the response from the Treasury API, validates it, and returns a result.
// An ExchangeRate object and nil error if the response is valid. Otherwise, it returns a nil object and an error.
func ProcessResponse(resp *http.Response, currencyName string) (*domain.ExchangeRate, error) {
	// Checks the response status code after a successful request
	if resp.StatusCode != http.StatusOK {
		log.Warn().Int("status_code", resp.StatusCode).Str("currency", currencyName).Msg("unexpected API response")
		return nil, ErrTreasuryAPIResponse
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	var data struct {
		Data []struct {
			Currency     string `json:"currency"`
			ExchangeRate string `json:"exchange_rate"`
			RecordDay    string `json:"record_calendar_day"`
			RecordMonth  string `json:"record_calendar_month"`
			RecordYear   string `json:"record_calendar_year"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, ErrDecodingResponse
	}

	if len(data.Data) == 0 {
		log.Error().Str("currency", currencyName).Msg("no data found in API response")
		return nil, ErrExchangeRateNotFound
	}

	// Parse and validate the date of record
	dayOfRecord, monthOfRecord, yearOfRecord := data.Data[0].RecordDay, data.Data[0].RecordMonth, data.Data[0].RecordYear
	dateOfRecord, err := ParseDateFromResponse(dayOfRecord, monthOfRecord, yearOfRecord)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParsingExchangeRateDateOfRecord, err)
	}

	rate, err := strconv.ParseFloat(data.Data[0].ExchangeRate, 64)
	if err != nil {
		return nil, ErrInvalidExchangeRate
	}

	exchangeRate, errs := domain.NewExchangeRate(data.Data[0].Currency, rate, dateOfRecord)

	// If there are errors, joins them into one and return
	if len(errs) > 0 {
		var errMessages []string
		for _, e := range errs {
			errMessages = append(errMessages, e.Error())
		}
		return nil, fmt.Errorf("validation errors: %s", strings.Join(errMessages, ", "))
	}

	return exchangeRate, nil
}

// ParseDateFromResponse takes day, month, and year of record as strings, validates them,
// and returns a parsed time.Time or an error if the values are invalid.
func ParseDateFromResponse(dayString, monthString, yearString string) (time.Time, error) {
	dayString = strings.TrimSpace(dayString)
	monthString = strings.TrimSpace(monthString)
	yearString = strings.TrimSpace(yearString)

	// Validates the day, month, and year parameters emptiness: must not be empty
	if yearString == "" || monthString == "" || dayString == "" {
		return time.Time{}, ErrEmptyField
	}

	day, err := strconv.Atoi(dayString)
	if err != nil || day < 1 || day > 31 {
		return time.Time{}, ErrInvalidDay
	}
	month, err := strconv.Atoi(monthString)
	if err != nil || month < 1 || month > 12 {
		return time.Time{}, ErrInvalidMonth
	}
	year, err := strconv.Atoi(yearString)
	if err != nil {
		return time.Time{}, ErrInvalidYear
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), nil
}
