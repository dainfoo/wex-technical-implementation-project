package domain

import (
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"
)

// This file contains the definition of the Transaction, in a struct, and the NewTransaction constructor function.
// It also includes helper functions for input validation and conversion.

// Transaction represents a financial transaction.
type Transaction struct {
	// ID is the unique identifier for the transaction.
	ID uuid.UUID
	// Description provides details about the transaction. It must not be empty and must not exceed 50 characters.
	Description string
	// Timestamp is the time when the transaction occurred, stored in UTC.
	Timestamp time.Time
	// AmountInUSD is the transaction amount in USD, rounded to two decimal places.
	AmountInUSD *big.Float
}

// NewTransaction creates a new Transaction instance with input validation.
func NewTransaction(description string, timestampString string, amount float64) (*Transaction, []error) {
	description = strings.TrimSpace(description)
	timestampString = strings.TrimSpace(timestampString)

	var errors []error = make([]error, 0, 5)

	// Aggregate errors from each validation function
	errors = append(errors, ValidateDescription(description)...)
	errors = append(errors, ValidateAmount(amount)...)
	timestamp, timestampErrors := ParseAndValidateTimestamp(timestampString)
	errors = append(errors, timestampErrors...)

	// Check if validation errors occur, returning them and stopping the transaction creation if any are found
	if len(errors) > 0 {
		return nil, errors
	}

	amountInUSD := new(big.Float).SetPrec(64).SetFloat64(RoundToTwoDecimalPlaces(amount))
	id := uuid.New()

	return &Transaction{
		ID:          id,
		Description: description,
		Timestamp:   timestamp.UTC(),
		AmountInUSD: amountInUSD,
	}, errors
}

// ValidateDescription validates the transaction description.
func ValidateDescription(description string) []error {
	var errors []error = make([]error, 0, 2)

	// Validate the description length: must not be empty
	if len(description) == 0 {
		errors = append(errors, ErrDescriptionEmpty)
	}
	// Validate the description length: must not exceed 50 characters
	if len(description) > 50 {
		errors = append(errors, ErrDescriptionTooLong)
	}

	return errors
}

// ValidateAmount validates the transaction amount.
func ValidateAmount(amount float64) []error {
	var errors []error = make([]error, 0, 1)

	// Validate the amount: must be positive
	if amount <= 0 {
		errors = append(errors, ErrInvalidAmount)
	}

	return errors
}

// ParseISO8601Timestamp validates if the provided timestamp string is in ISO 8601 format
// and converts it to a time.Time instance.
func ParseISO8601Timestamp(timestampString string) (time.Time, error) {
	timestampString = strings.TrimSpace(timestampString)
	const layout = "2006-01-02T15:04:05Z07:00" // ISO 8601 layout

	// Validate the timestamp string: must not be empty
	if timestampString == "" {
		return time.Time{}, ErrTimestampEmpty
	}

	// Attempt to parse the timestamp string
	parsedTimestamp, err := time.Parse(layout, timestampString)
	if err != nil {
		return time.Time{}, ErrInvalidTimestampFormat
	}

	return parsedTimestamp, nil
}

// ParseAndValidateTimestamp checks if the provided timestamp string is not empty, parses it,
// and ensures that the timestamp is not in the future.
func ParseAndValidateTimestamp(timestampString string) (time.Time, []error) {
	var errors []error = make([]error, 0, 1)

	// Validate the timestamp string: must not be empty
	if timestampString == "" {
		errors = append(errors, ErrTimestampEmpty)
		return time.Time{}, errors
	}

	timestamp, err := ParseISO8601Timestamp(timestampString)
	if err != nil {
		errors = append(errors, err)
		return time.Time{}, errors
	}

	// Validate the parsed timestamp: must not be in the future
	if timestamp.After(time.Now()) {
		errors = append(errors, ErrInvalidTimestamp)
		return time.Time{}, errors
	}

	return timestamp, errors
}

// RoundToTwoDecimalPlaces rounds a float64 to two decimal places.
func RoundToTwoDecimalPlaces(amount float64) float64 {
	return math.Round(amount*100) / 100
}
