package domain

import "errors"

// This file defines error variables related to transaction validation in the domain layer.

var (
	// ErrDescriptionEmpty is returned when the transaction description is empty.
	ErrDescriptionEmpty = errors.New("transaction description is required; it cannot be empty")

	// ErrDescriptionTooLong is returned when the transaction description exceeds the allowed character limit.
	ErrDescriptionTooLong = errors.New("transaction description is invalid; it must not exceed 50 characters")

	// ErrInvalidTimestamp is returned when the transaction timestamp is invalid.
	ErrInvalidTimestamp = errors.New("transaction timestamp is invalid; it cannot be in the future")

	// ErrInvalidAmountInUSD is returned when the transaction amount in USD is invalid.
	ErrInvalidAmountInUSD = errors.New("transaction amount in USD is invalid; it must be greater than 0")
)
