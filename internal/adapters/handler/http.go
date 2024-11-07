package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/dainfoo/wex-technical-implementation-project/internal/core/domain"
	"github.com/dainfoo/wex-technical-implementation-project/internal/core/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// This file contains the HTTP handler for transactions.

// TransactionHandler holds the resources needed to handle HTTP requests for transactions.
type TransactionHandler struct {
	transactionService services.TransactionService
}

// TransactionDTO represents the data transfer object for transactions.
type TransactionDTO struct {
	ID                     string  `json:"id"`
	Description            string  `json:"description"`
	Timestamp              string  `json:"timestamp"`
	AmountInUSD            float64 `json:"amount_in_usd"`
	ExchangeRateUsed       float64 `json:"exchange_rate_used"`
	AmountInTargetCurrency float64 `json:"amount_in_target_currency"`
}

// SuccessResponse wraps successful responses.
type SuccessResponse struct {
	Data interface{} `json:"data"`
}

// ErrorResponse wraps error responses.
type ErrorResponse struct {
	Error string `json:"error"`
}

// NewTransactionHandler creates a new handler with injected services.
func NewTransactionHandler(transactionService services.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

// Routes sets up the Chi router with the necessary routes.
func (th *TransactionHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(httprate.LimitByIP(100, time.Minute))
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			w.Header().Set("X-Frame-Options", "DENY")
			next.ServeHTTP(w, r)
		})
	})

	r.Post("/transactions", th.SaveTransaction)
	r.Get("/transactions/{id}/{currency}", th.FindTransactionWithCurrencyConversion)
	r.Get("/health", th.HealthCheck)

	return r
}

// SaveTransaction handles the POST request to save a new transaction.
func (th *TransactionHandler) SaveTransaction(w http.ResponseWriter, r *http.Request) {
	var data TransactionDTO = TransactionDTO{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Warn().Err(err).Msg("invalid request payload")
		WriteErrorResponse(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	transaction, validationErrors := th.ValidateAndCreateTransaction(data)
	if len(validationErrors) > 0 {
		validationErrorsAsStrings := make([]string, len(validationErrors))
		for i, err := range validationErrors {
			validationErrorsAsStrings[i] = err.Error()
		}
		joinedErrors := strings.Join(validationErrorsAsStrings, ", ")
		log.Warn().Errs("validation_errors", validationErrors).Str("transaction_id",
			data.ID).Msg("transaction validation failed")
		WriteErrorResponse(w, http.StatusBadRequest, "validation errors: "+joinedErrors)
		return
	}

	if err := th.transactionService.SaveTransaction(*transaction); err != nil {
		log.Error().Err(err).Msg("failed to save the transaction")
		WriteErrorResponse(w, http.StatusInternalServerError, "failed to save the transaction")
		return
	}

	WriteSuccessResponse(w, map[string]string{"id": transaction.ID.String()}, http.StatusCreated)
}

// FindTransactionWithCurrencyConversion handles the GET request to find and return a transaction
// converted to a target currency.
func (th *TransactionHandler) FindTransactionWithCurrencyConversion(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	id, err := uuid.Parse(idString)
	if err != nil {
		log.Warn().Err(err).Str("id", idString).Msg("invalid transaction ID format")
		WriteErrorResponse(w, http.StatusBadRequest, "invalid transaction ID format")
		return
	}
	currencyName := chi.URLParam(r, "currency")
	if currencyName == "" {
		log.Warn().Msg("currency not provided")
		WriteErrorResponse(w, http.StatusBadRequest, "currency not provided")
		return
	}
	transaction, exchangeRate, err := th.transactionService.FindTransactionAndExchangeRateFromCurrency(id, currencyName)
	if err != nil {
		log.Warn().Err(err).Msg("transaction not found or cannot be converted to the target currency")
		WriteErrorResponse(w, http.StatusNotFound, "the purchase cannot be converted to the target currency")
		return
	}
	transactionAmountInUSD, _ := transaction.AmountInUSD.Float64()
	exchangeRateUsed, _ := exchangeRate.Rate.Float64()
	transactionAmountInUSD = domain.RoundToTwoDecimalPlaces(transactionAmountInUSD)
	exchangeRateUsed = domain.RoundToTwoDecimalPlaces(exchangeRateUsed)

	transactionDTO := TransactionDTO{
		ID:                     transaction.ID.String(),
		Description:            transaction.Description,
		Timestamp:              transaction.Timestamp.Format(time.DateTime),
		AmountInUSD:            transactionAmountInUSD,
		ExchangeRateUsed:       exchangeRateUsed,
		AmountInTargetCurrency: domain.RoundToTwoDecimalPlaces(transactionAmountInUSD * exchangeRateUsed),
	}

	WriteSuccessResponse(w, transactionDTO, http.StatusOK)
}

// HealthCheck handles the GET request to check the health of the server.
func (th *TransactionHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now(),
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(health); err != nil {
		log.Error().Err(err).Msg("failed to encode response")
		WriteErrorResponse(w, http.StatusInternalServerError, "internal server error")
	}
}

// ValidateAndCreateTransaction validates and creates a new transaction from the provided request data.
func (th *TransactionHandler) ValidateAndCreateTransaction(data TransactionDTO) (*domain.Transaction, []error) {
	timestamp, errs := ParseAndValidateTimestamp(data.Timestamp)
	if len(errs) > 0 {
		return nil, errs
	}
	return domain.NewTransaction(data.Description, timestamp, data.AmountInUSD)
}

// StartServer starts the HTTP server on the provided port.
func (th *TransactionHandler) StartServer(port string) {
	router := th.Routes()
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  0 * time.Second,
	}

	// Start the server in a goroutine
	go func() {
		log.Info().Str("port", port).Msg("starting the server")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("server failed to start")
		}
	}()

	// Check if the server is ready
	go func() {
		for {
			resp, err := http.Get("http://localhost" + server.Addr + "/health")
			if err == nil && resp.StatusCode == http.StatusOK {
				log.Printf("server started successfully on port %s", port)
				if err := resp.Body.Close(); err != nil {
					log.Warn().Err(err).Msg("error closing response body")
				}
				break
			}
			// Wait for 500 milliseconds before checking again
			time.Sleep(500 * time.Millisecond)
		}
	}()

	th.ShutdownServer(server)
}

// ShutdownServer gracefully shuts down the server when an interrupt signal is received.
func (th *TransactionHandler) ShutdownServer(server *http.Server) {
	// Gracefully shutdown the server
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Info().Msg("shutting down the server")
	// Wait for 5 seconds before shutting down the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("server forced to shutdown")
	}
	log.Info().Msg("server exited")
}

// ParseAndValidateTimestamp checks if the provided timestamp string is not empty, parses it,
// and ensures that the timestamp is not in the future.
func ParseAndValidateTimestamp(timestampString string) (time.Time, []error) {
	var errs []error = make([]error, 0, 1)

	// Validate the timestamp string: must not be empty
	if timestampString == "" {
		errs = append(errs, ErrTimestampEmpty)
		return time.Time{}, errs
	}

	timestamp, err := ParseISO8601Timestamp(timestampString)
	if err != nil {
		errs = append(errs, err)
		return time.Time{}, errs
	}

	// Validate the parsed timestamp: must not be in the future
	if timestamp.After(time.Now()) {
		errs = append(errs, ErrInvalidTimestamp)
		return time.Time{}, errs
	}

	return timestamp, errs
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

// WriteSuccessResponse writes a success response with the provided data and ensures the status code is set only once.
func WriteSuccessResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	// Defaults to http.StatusOK if none is provided (0)
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	// Write the response with the data wrapped in the "data" field
	response := SuccessResponse{Data: data}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error().Err(err).Msg("failed to encode response")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// WriteErrorResponse writes an error response with the provided status code and message.
func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Wrap the error message in the "error" field
	response := ErrorResponse{Error: message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error().Err(err).Msg("failed to encode response")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
