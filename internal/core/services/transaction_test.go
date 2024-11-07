package services_test

import (
	"os"
	"testing"
	"time"

	"github.com/dainfoo/wex-technical-implementation-project/internal/adapters/client"
	"github.com/dainfoo/wex-technical-implementation-project/internal/adapters/repository"
	"github.com/dainfoo/wex-technical-implementation-project/internal/core/domain"
	"github.com/dainfoo/wex-technical-implementation-project/internal/core/ports"
	"github.com/dainfoo/wex-technical-implementation-project/internal/core/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// This file contains a test suite for the TransactionService adapter.
// It uses Testify for assertions and mocking, and runs the tests in parallel.

// TransactionServiceIntegrationTestSuite represents the test suite.
type TransactionServiceIntegrationTestSuite struct {
	suite.Suite
	transactionRepo ports.TransactionRepository
	exchangeAdapter *client.MockTreasuryExchangeRateAdapter
	service         *services.TransactionService
}

// SetupTest initializes the test suite.
func (suite *TransactionServiceIntegrationTestSuite) SetupTest() {
	testDatabasePath := "service_test.db"
	boltDBRepo, err := repository.NewTransactionRepositoryBoltDB(testDatabasePath, "transactions")
	suite.NoError(err)

	mockAdapter := new(client.MockTreasuryExchangeRateAdapter)

	suite.service = services.NewTransactionService(boltDBRepo, mockAdapter)
	suite.transactionRepo = boltDBRepo
	suite.exchangeAdapter = mockAdapter
	// Clean up the database after the test suite finishes
	suite.T().Cleanup(func() {
		if repo, ok := suite.transactionRepo.(*repository.TransactionRepositoryBoltDB); ok {
			err := repo.GetBoltDB().Close()
			require.NoError(suite.T(), err, "failed to close BoltDB")

			err = os.Remove(testDatabasePath)
			require.NoError(suite.T(), err, "failed to delete test database file")
		}
	})
}

// TestFindTransactionAndExchangeRate tests the FindTransactionAndExchangeRate method of the TransactionService.
func (suite *TransactionServiceIntegrationTestSuite) TestFindTransactionAndExchangeRate() {
	// Expected results
	successTransaction, err := domain.NewTransaction("giberish", time.Now(), 25.7)
	// Stops the test if the expected results are not as expected (probably the business logic changed)
	require.Empty(suite.T(), err)
	successExchangeRate, err := domain.NewExchangeRate("Real", 5.434, time.Date(2024, 9, 30, 0, 0, 0, 0, time.UTC))
	require.Empty(suite.T(), err)

	tests := []struct {
		name                string
		transactionID       uuid.UUID
		setupTransaction    *domain.Transaction
		currencyName        string
		mockRate            *domain.ExchangeRate
		mockRateErr         error
		expectedTransaction *domain.Transaction
		expectedRate        *domain.ExchangeRate
		expectedErr         error
	}{
		{
			name:                "Success",
			transactionID:       successTransaction.ID,
			setupTransaction:    successTransaction,
			currencyName:        "Real",
			mockRate:            successExchangeRate,
			expectedTransaction: successTransaction,
			expectedRate:        successExchangeRate,
			expectedErr:         nil,
		},
		{
			name:                "Transaction Not Found",
			transactionID:       uuid.New(),
			currencyName:        "Real",
			expectedTransaction: nil,
			expectedRate:        nil,
			expectedErr:         repository.ErrTransactionNotFound,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Tests the setup of the test case
			if tt.setupTransaction != nil {
				err := suite.transactionRepo.SaveTransaction(*tt.setupTransaction)
				suite.NoError(err)
			}
			suite.exchangeAdapter.On("GetExchangeRates", tt.currencyName).
				Return([]*domain.ExchangeRate{tt.mockRate}, tt.mockRateErr)

			foundTransaction, exchangeRate, err := suite.service.FindTransactionAndExchangeRateFromCurrency(tt.transactionID, tt.currencyName)

			if tt.expectedErr != nil {
				assert.Error(suite.T(), err)
				assert.Equal(suite.T(), tt.expectedErr, err)
			} else {
				suite.NoError(err)
				assert.Equal(suite.T(), tt.expectedTransaction.ID.String(), foundTransaction.ID.String())
				assert.Equal(suite.T(), tt.expectedTransaction.Description, foundTransaction.Description)
				assert.Equal(suite.T(), tt.expectedTransaction.Timestamp, foundTransaction.Timestamp)
				assert.Equal(suite.T(), tt.expectedTransaction.AmountInUSD.Cmp(foundTransaction.AmountInUSD), 0)
				assert.Equal(suite.T(), tt.expectedRate.CurrencyName, exchangeRate.CurrencyName)
				assert.Equal(suite.T(), tt.expectedRate.Rate.Cmp(exchangeRate.Rate), 0)
			}

			// Reset mock expectations for the next test case
			suite.exchangeAdapter.ExpectedCalls = nil
		})
	}
}

// TestTransactionServiceIntegrationTestSuite initializes the test suite.
func TestTransactionServiceIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionServiceIntegrationTestSuite))
}
