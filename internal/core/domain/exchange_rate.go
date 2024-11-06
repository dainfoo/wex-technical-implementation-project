package domain

import (
	"math/big"
	"strings"
	"time"
)

// This file contains the ExchangeRate struct, its constructor and validation functions.

// ExchangeRate represents an exchange rate.
type ExchangeRate struct {
	CurrencyName string
	Rate         *big.Float
	DateOfRecord time.Time
}

// NewExchangeRate creates a new ExchangeRate instance with input validation.
func NewExchangeRate(currencyName string, rate float64, dateOfRecord time.Time) (*ExchangeRate, []error) {
	currencyName = strings.TrimSpace(currencyName)

	// Validate the inputs before constructing the object
	if errs := ValidateExchangeRate(currencyName, rate, dateOfRecord); len(errs) > 0 {
		return nil, errs
	}

	return &ExchangeRate{
		CurrencyName: currencyName,
		Rate:         new(big.Float).SetPrec(64).SetFloat64(rate),
		DateOfRecord: dateOfRecord,
	}, nil
}

// ValidateExchangeRate validates the currency name, rate and date of record for the ExchangeRate struct.
func ValidateExchangeRate(currencyName string, rate float64, dateOfRecord time.Time) []error {
	var errors []error = make([]error, 0, 3)

	// Validate the currency name length: must not be empty
	if currencyName == "" {
		errors = append(errors, ErrCurrencyNameEmpty)
	}

	// Validate the rate: must be positive
	if rate <= 0 {
		errors = append(errors, ErrInvalidExchangeRate)
	}

	// Validate the date of record: cannot be in the future
	if dateOfRecord.After(time.Now()) {
		errors = append(errors, ErrInvalidDateOfRecord)
	}

	return errors
}
