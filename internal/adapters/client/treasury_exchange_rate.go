package client

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dainfoo/wex-technical-implementation-project/internal/core/domain"
	"github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
)

// This file contains the implementation of the ExchangeRateService	interface using the Treasury API.

// Activate the jsoniter library to decode the Treasury API response.
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// TreasuryExchangeRateAdapter interface defines the behavior for exchange rates fetching.
// It allows flexibility to change the implementation of the Treasury API client for testing purposes.
type TreasuryExchangeRateAdapter interface {
	GetExchangeRates(currencyName string) ([]*domain.ExchangeRate, error)
}

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

// ConcreteTreasuryExchangeRateAdapter is the real implementation of TreasuryExchangeRateAdapter interface.
type ConcreteTreasuryExchangeRateAdapter struct {
	client      HTTPClient
	apiEndpoint string
}

// NewConcreteTreasuryExchangeRateAdapter creates a new ConcreteTreasuryExchangeRateAdapter with the given HTTPClient.
func NewConcreteTreasuryExchangeRateAdapter(client HTTPClient) *ConcreteTreasuryExchangeRateAdapter {
	return &ConcreteTreasuryExchangeRateAdapter{
		client:      client,
		apiEndpoint: treasuryAPIEndpoint,
	}
}

// GetExchangeRates retrieves all the exchange rates for a currency with input and response validations.
func (a *ConcreteTreasuryExchangeRateAdapter) GetExchangeRates(currencyName string) ([]*domain.ExchangeRate, error) {
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
			log.Warn().Err(err).Int("attempt", attempt+1).Msg("retrying request due to transient network issue")
			time.Sleep(retryDelay)
		}
	}
	if err != nil {
		log.Error().Err(err).Msg("error fetching exchange rates from Treasury API")
		return nil, ErrNetworkIssue
	}

	return ProcessResponse(resp, currencyName)
}

// buildRequestURL constructs the URL for the Treasury API request.
func buildRequestURL(a *ConcreteTreasuryExchangeRateAdapter, currencyName string) string {
	return fmt.Sprintf("%s?&sort=-record_date&format=json&page[number]=1&page[size]=1000"+
		"&fields=currency,exchange_rate,record_date,record_calendar_day,record_calendar_month,record_calendar_year"+
		"&filter=currency:eq:%s", a.apiEndpoint, url.QueryEscape(currencyName))
}

// ProcessResponse reads the response from the Treasury API, validates it, and returns a result.
// An ExchangeRate slice and nil error if the response is valid. Otherwise, it returns a nil object and an error.
func ProcessResponse(resp *http.Response, currencyName string) ([]*domain.ExchangeRate, error) {
	// Checks the response status code after a successful request
	if resp.StatusCode != http.StatusOK {
		log.Error().Int("status_code", resp.StatusCode).Str("currency", currencyName).Msg("unexpected API response")
		return nil, ErrTreasuryAPIResponse
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error().Err(err).Msg("error closing response body")
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
		log.Error().Err(err).Msg("error decoding API response")
		return nil, ErrDecodingResponse
	}

	if len(data.Data) == 0 {
		log.Error().Str("currency", currencyName).Msg("no data found in API response")
		return nil, ErrExchangeRateNotFound
	}

	var exchangeRates []*domain.ExchangeRate
	// Loops through all the data and parse each exchange rate
	for _, item := range data.Data {
		dayOfRecord, monthOfRecord, yearOfRecord := item.RecordDay, item.RecordMonth, item.RecordYear
		dateOfRecord, err := ParseDateFromResponse(dayOfRecord, monthOfRecord, yearOfRecord)
		if err != nil {
			return nil, fmt.Errorf("error parsing exchange rate date of record (day: %s, month: %s, year: %s): %w",
				dayOfRecord, monthOfRecord, yearOfRecord, ErrParsingExchangeRateDateOfRecord)
		}

		rate, err := strconv.ParseFloat(item.ExchangeRate, 64)
		if err != nil {
			return nil, ErrInvalidExchangeRate
		}

		exchangeRate, errs := domain.NewExchangeRate(item.Currency, rate, dateOfRecord)

		// If there are errors, join them into one and return
		if len(errs) > 0 {
			var errMessages []string
			for _, e := range errs {
				errMessages = append(errMessages, e.Error())
			}
			return nil, fmt.Errorf("validation errors: %s", strings.Join(errMessages, ", "))
		}

		exchangeRates = append(exchangeRates, exchangeRate)
	}

	return exchangeRates, nil
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
