package domain

import (
	"math/big"
	"time"
)

// This file contains the ExchangeRate struct, its constructor and validation functions.

// ExchangeRate represents an exchange rate.
type ExchangeRate struct {
	// Currency is the currency name.
	Currency string
	// Rate is the exchange rate.
	Rate *big.Float
	// Date is the date of the exchange rate record.
	Date time.Time
}

// NewExchangeRate creates a new ExchangeRate instance with input validation.
func NewExchangeRate(currency string, rate float64, date time.Time) (*ExchangeRate, []error) {
	// Validate the inputs before constructing the object
	if errs := ValidateExchangeRate(currency, rate, date); len(errs) > 0 {
		return nil, errs
	}

	return &ExchangeRate{
		Currency: currency,
		Rate:     new(big.Float).SetPrec(64).SetFloat64(rate),
		Date:     date,
	}, nil
}

// ValidateExchangeRate validates the currency, rate and date for the ExchangeRate struct.
func ValidateExchangeRate(currency string, rate float64, date time.Time) []error {
	var errors []error = make([]error, 0, 3)

	// Validate the currency name length: must not be empty
	if currency == "" {
		errors = append(errors, ErrCurrencyEmpty)
	}

	// Validate the rate: must be positive
	if rate <= 0 {
		errors = append(errors, ErrInvalidExchangeRate)
	}

	// Validate the date: cannot be in the future
	if date.After(time.Now()) {
		errors = append(errors, ErrInvalidDate)
	}

	return errors
}
