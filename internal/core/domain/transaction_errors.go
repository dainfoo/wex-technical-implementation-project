package domain

import "errors"

// This file defines error variables related to transaction validation in the domain layer.

var (
	// ErrDescriptionEmpty is returned when the description is empty.
	ErrDescriptionEmpty = errors.New("transaction description is required; it cannot be empty")

	// ErrDescriptionTooLong is returned when the description exceeds the allowed character limit.
	ErrDescriptionTooLong = errors.New("transaction description must not exceed 50 characters")

	// ErrInvalidAmount is returned when the transaction amount is invalid.
	ErrInvalidAmount = errors.New("transaction amount must be a positive value")

	// ErrTimestampEmpty is returned when the timestamp is empty.
	ErrTimestampEmpty = errors.New("transaction timestamp is required; it cannot be empty")

	// ErrInvalidTimestampFormat is returned when the timestamp format is invalid.
	ErrInvalidTimestampFormat = errors.New("transaction timestamp format must be in ISO 8601 standard")

	// ErrInvalidTimestamp is returned when the timestamp is in the future.
	ErrInvalidTimestamp = errors.New("transaction timestamp cannot be in the future")
)
