package domain

import (
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"
)

// This file contains the Transaction struct, its constructor and validation functions.

// Transaction represents a financial transaction.
type Transaction struct {
	// ID is the unique identifier for the transaction.
	ID uuid.UUID `json:"id"`
	// Description provides details about the transaction. It cannot be empty and must not exceed 50 characters.
	Description string `json:"description"`
	// Timestamp is the time when the transaction occurred, stored in UTC.
	Timestamp time.Time `json:"timestamp"`
	// AmountInUSD is the transaction amount in USD, rounded to two decimal places.
	AmountInUSD *big.Float `json:"amount_in_usd"`
}

// NewTransaction creates a new Transaction instance with input validation.
func NewTransaction(description string, timestamp time.Time, amountInUSD float64) (*Transaction, []error) {
	description = strings.TrimSpace(description)

	// Validate the inputs before constructing the object and stop the transaction creation if any errors are found
	if errs := ValidateTransaction(description, timestamp, amountInUSD); len(errs) > 0 {
		return nil, errs
	}

	amountInUSDBigFloat := new(big.Float).SetPrec(64).SetFloat64(RoundToTwoDecimalPlaces(amountInUSD))
	id := uuid.New()

	return &Transaction{
		ID:          id,
		Description: description,
		Timestamp:   timestamp.UTC(),
		AmountInUSD: amountInUSDBigFloat,
	}, nil
}

// ValidateTransaction validates the description, timestamp and the amount in USD for the Transaction struct.
func ValidateTransaction(description string, timestamp time.Time, amountInUSD float64) []error {
	var errors []error = make([]error, 0, 5)

	// Aggregate the validation errors
	errors = append(errors, ValidateDescription(description)...)
	errors = append(errors, ValidateAmountInUSD(amountInUSD)...)

	// Validate the timestamp: cannot be in the future
	if timestamp.After(time.Now()) {
		errors = append(errors, ErrInvalidTimestamp)
	}

	return errors
}

// ValidateDescription validates the transaction description.
func ValidateDescription(description string) []error {
	var errors []error = make([]error, 0, 1)

	// Validate the description emptiness: must not be empty
	if len(description) == 0 {
		errors = append(errors, ErrDescriptionEmpty)
		return errors
	}
	// Validate the description length: must not exceed 50 characters
	if len(description) > 50 {
		errors = append(errors, ErrDescriptionTooLong)
	}

	return errors
}

// ValidateAmountInUSD validates the transaction amount in USD.
func ValidateAmountInUSD(amountInUSD float64) []error {
	var errors []error = make([]error, 0, 1)

	// Validate the amount in USD: must be positive
	if amountInUSD <= 0 {
		errors = append(errors, ErrInvalidAmountInUSD)
	}

	return errors
}

// RoundToTwoDecimalPlaces rounds a float64 to two decimal places.
func RoundToTwoDecimalPlaces(value float64) float64 {
	return math.Round(value*100) / 100
}
