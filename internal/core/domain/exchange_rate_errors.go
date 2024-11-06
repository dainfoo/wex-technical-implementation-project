package domain

import "errors"

// This file defines error variables related to exchange rate validation in the domain layer.

var (
	// ErrCurrencyEmpty is returned when the currency is empty.
	ErrCurrencyEmpty = errors.New("currency name is required; it cannot be empty")

	// ErrInvalidExchangeRate is returned when the exchange rate is invalid.
	ErrInvalidExchangeRate = errors.New("exchange rate is invalid; it must be greater than 0")

	// ErrInvalidDate is returned when the date is invalid.
	ErrInvalidDate = errors.New("date is invalid; it cannot be in the future")
)
