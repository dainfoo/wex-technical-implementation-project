package client

import "errors"

// This file defines error variables related to the ExchangeRateService port implementation using the Treasury API.

var (
	// ErrEmptyField is returned when one of the expected fields (like day, month, year) is empty.
	ErrEmptyField = errors.New("day, month and/or year must not be empty")

	// ErrInvalidDay is returned when the provided day is invalid or cannot be parsed.
	ErrInvalidDay = errors.New("invalid day")

	// ErrInvalidMonth is returned when the provided month is invalid or cannot be parsed.
	ErrInvalidMonth = errors.New("invalid month")

	// ErrInvalidYear is returned when the provided year is invalid or cannot be parsed.
	ErrInvalidYear = errors.New("invalid year")

	// ErrParsingExchangeRateDateOfRecord is returned when the exchange rate date of record cannot be parsed.
	ErrParsingExchangeRateDateOfRecord = errors.New("error parsing exchange rate date of record")

	// ErrInvalidExchangeRate is returned when the exchange rate value is invalid and cannot be parsed into a big.Float.
	ErrInvalidExchangeRate = errors.New("invalid exchange rate value")

	// ErrExchangeRateNotFound is returned when no exchange rate is found for the requested currency.
	ErrExchangeRateNotFound = errors.New("no exchange rate found for the provided currency")

	// ErrDecodingResponse is returned when the API response cannot be decoded into the expected format.
	ErrDecodingResponse = errors.New("error decoding response from Treasury API")

	// ErrNetworkIssue is returned when there is a network-related issue when fetching data from the Treasury API.
	ErrNetworkIssue = errors.New("network issue while fetching exchange rate data")

	// ErrTreasuryAPIResponse is returned when the Treasury API returns an error.
	ErrTreasuryAPIResponse = errors.New("error from the Treasury API")
)
