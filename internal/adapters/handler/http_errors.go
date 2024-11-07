package handler

import "errors"

// This file defines error variables related to the HTTP handler.

var (
	// ErrTimestampEmpty is returned when the timestamp is empty.
	ErrTimestampEmpty = errors.New("transaction timestamp is required; it cannot be empty")

	// ErrInvalidTimestampFormat is returned when the timestamp format is invalid.
	ErrInvalidTimestampFormat = errors.New("transaction timestamp format must be in ISO 8601 standard")

	// ErrInvalidTimestamp is returned when the timestamp is in the future.
	ErrInvalidTimestamp = errors.New("transaction timestamp cannot be in the future")
)
