package domain

import "errors"

// This file defines error variables related to exchange rate validation in the domain layer.

var (
	// ErrCurrencyNameEmpty is returned when the currency name is empty.
	ErrCurrencyNameEmpty = errors.New("exchange rate currency name is required; it cannot be empty")

	// ErrInvalidExchangeRate is returned when the exchange rate is invalid.
	ErrInvalidExchangeRate = errors.New("exchange rate is invalid; it must be greater than 0")

	// ErrInvalidDateOfRecord is returned when the date of record is invalid.
	ErrInvalidDateOfRecord = errors.New("exchange rate date of record is invalid; it cannot be in the future")
)
