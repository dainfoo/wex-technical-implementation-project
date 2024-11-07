package repository_test

import (
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/dainfoo/wex-technical-implementation-project/internal/adapters/repository"
	"github.com/dainfoo/wex-technical-implementation-project/internal/core/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.etcd.io/bbolt"
)

// This file contains tests for the BoltDB implementation of the TransactionRepository interface.
// It uses Testify for assertions and runs the tests in parallel.

// TestTransactionBoltDBRepository tests the BoltDB implementation of the TransactionRepository interface.
// It tests the following scenarios:
//
// 1. Repository Initialization.
// 2. Save And Find A Transaction.
// 3. Retrieve Non-Existent Transaction.
// 4. Bucket Not Found Error.
// 5. Persistence Across Sessions.
// 6. Heavy Write Scenario.
// 7. Heavy Read/Write Scenario.
func TestTransactionBoltDBRepository(t *testing.T) {
	// Reusable test Transaction
	testTransaction, err := domain.NewTransaction("giberish", time.Now(), 100.50)
	// Stops the test if the expected results are not as expected (probably the business logic changed)
	require.Empty(t, err)

	// Create a temporary BoltDB database file for testing
	tempDBPath := "testdata/transaction_test.db"

	// Ensure cleanup happens after all tests
	t.Cleanup(func() {
		err := os.RemoveAll("testdata")
		require.NoError(t, err, "failed to clean up test data directory")
	})

	t.Run("Repository Initialization", func(t *testing.T) {
		t.Parallel()

		// Create a unique bucket name for this test
		bucketName := "transactions_" + uuid.New().String()

		repo, err := repository.NewTransactionRepositoryBoltDB(tempDBPath, bucketName)
		require.NoError(t, err)
		err = repo.Close()
		require.NoError(t, err, "failed to close the repository")
	})

	t.Run("Save And Find A Transaction", func(t *testing.T) {
		t.Parallel()

		// Create a unique bucket name for this test
		bucketName := "transactions_" + uuid.New().String()

		repo, err := repository.NewTransactionRepositoryBoltDB(tempDBPath, bucketName)
		require.NoError(t, err)
		defer func() {
			err := repo.Close()
			require.NoError(t, err, "failed to close the repository")
		}()

		err = repo.SaveTransaction(*testTransaction)
		require.NoError(t, err)

		retrievedTransaction, err := repo.FindTransaction(testTransaction.ID)
		require.NoError(t, err)
		require.NotNil(t, retrievedTransaction)
		assert.Equal(t, testTransaction.ID, retrievedTransaction.ID)
		assert.Equal(t, testTransaction.AmountInUSD, retrievedTransaction.AmountInUSD)
		assert.Equal(t, testTransaction.Timestamp.UTC(), retrievedTransaction.Timestamp.UTC())
	})

	t.Run("Retrieve Non-Existent Transaction", func(t *testing.T) {
		t.Parallel()

		// Create a unique bucket name for this test
		bucketName := "transactions_" + uuid.New().String()

		repo, err := repository.NewTransactionRepositoryBoltDB(tempDBPath, bucketName)
		require.NoError(t, err)
		defer func() {
			err := repo.Close()
			require.NoError(t, err, "failed to close the repository")
		}()

		_, err = repo.FindTransaction(uuid.New())
		assert.ErrorIs(t, err, repository.ErrTransactionNotFound)
	})

	t.Run("Bucket Not Found Error", func(t *testing.T) {
		t.Parallel()

		// Create a unique bucket name for this test
		bucketName := "transactions_" + uuid.New().String()

		repo, err := repository.NewTransactionRepositoryBoltDB(tempDBPath, bucketName)
		require.NoError(t, err)
		defer func() {
			err := repo.Close()
			require.NoError(t, err, "failed to close the repository")
		}()

		boltDB := repo.GetBoltDB()
		err = boltDB.Update(func(tx *bbolt.Tx) error {
			// Delete the bucket to simulate the missing bucket scenario
			return tx.DeleteBucket([]byte(bucketName))
		})
		require.NoError(t, err)

		err = repo.SaveTransaction(*testTransaction)
		assert.ErrorIs(t, err, repository.ErrBucketNotFound)
	})

	t.Run("Persistence Across Sessions", func(t *testing.T) {
		t.Parallel()

		// Create a unique bucket name for this test
		bucketName := "transactions_" + uuid.New().String()

		repoFirstSession, err := repository.NewTransactionRepositoryBoltDB(tempDBPath, bucketName)
		require.NoError(t, err)

		err = repoFirstSession.SaveTransaction(*testTransaction)
		require.NoError(t, err)
		err = repoFirstSession.Close()
		require.NoError(t, err, "failed to close the first repository")

		// Second session: reopen and retrieve
		repoSecondSession, err := repository.NewTransactionRepositoryBoltDB(tempDBPath, bucketName)
		require.NoError(t, err)
		defer func() {
			err := repoSecondSession.Close()
			require.NoError(t, err, "failed to close the second repository")
		}()

		retrievedTransaction, err := repoSecondSession.FindTransaction(testTransaction.ID)
		require.NoError(t, err)
		assert.Equal(t, testTransaction.AmountInUSD, retrievedTransaction.AmountInUSD)
	})

	t.Run("Heavy Write Scenario", func(t *testing.T) {
		t.Parallel()

		// Create a unique bucket name for this test
		bucketName := "transactions_" + uuid.New().String()

		repo, err := repository.NewTransactionRepositoryBoltDB(tempDBPath, bucketName)
		require.NoError(t, err)
		defer func() {
			err := repo.Close()
			require.NoError(t, err, "failed to close the repository")
		}()

		// Simulate multiple concurrent writes

		var wg sync.WaitGroup
		iterations := 50000
		transactions := make([]domain.Transaction, iterations)
		writeErrChan := make(chan error, iterations)

		// Write transactions concurrently
		for i := 0; i < iterations; i++ {
			// Creates a test transaction
			transaction, err := domain.NewTransaction("giberish", time.Now(), float64(i)+(rand.Float64()*100))
			// Stops the test if the expected results are not as expected (probably the business logic changed)
			require.Empty(t, err)

			transactions[i] = *transaction

			wg.Add(1)
			go func(transaction domain.Transaction) {
				defer wg.Done()
				err := repo.SaveTransaction(transaction)
				if err != nil {
					writeErrChan <- err
				}
			}(*transaction)
		}

		// Waits for all goroutines to finish
		wg.Wait()
		// Closes the channel after all writes are done
		close(writeErrChan)

		// Process errors
		for err := range writeErrChan {
			assert.NoError(t, err, "error saving transaction")
		}

		// Verify if all transactions were saved
		for _, transaction := range transactions {
			retrievedTransaction, err := repo.FindTransaction(transaction.ID)
			require.NoError(t, err)
			assert.Equal(t, 0, transaction.AmountInUSD.Cmp(retrievedTransaction.AmountInUSD))
		}
	})

	t.Run("Heavy Read/Write Scenario", func(t *testing.T) {
		t.Parallel()

		// Creates a unique bucket name for this test
		bucketName := "transactions_" + uuid.New().String()

		repo, err := repository.NewTransactionRepositoryBoltDB(tempDBPath, bucketName)
		require.NoError(t, err)
		defer func() {
			err := repo.Close()
			require.NoError(t, err, "failed to close the repository")
		}()

		// Simulate multiple concurrent read/writes

		var wg sync.WaitGroup
		iterations := 25000
		writeErrChan := make(chan error, iterations)

		for i := 0; i < iterations; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()

				// Creates a test transaction
				transaction, err := domain.NewTransaction("giberish", time.Now(), float64(i)+(rand.Float64()*100))
				// Stops the test if the expected results are not as expected (probably the business logic changed)
				require.Empty(t, err)

				// Random write
				if err := repo.SaveTransaction(*transaction); err != nil {
					writeErrChan <- err
					return
				}

				// Random read
				if i%2 == 0 {
					if _, err := repo.FindTransaction(transaction.ID); err != nil {
						writeErrChan <- err
					}
				}
			}(i)
		}

		// Waits for all goroutines to finish
		wg.Wait()
		// Closes the channel after all writes are done
		close(writeErrChan)

		// Process errors
		for err := range writeErrChan {
			assert.NoError(t, err, "error during heavy read/write load operation")
		}
	})
}

// TestValidateTransactionRepositoryBoltDB tests the ValidateTransactionRepositoryBoltDB function.
// It tests the following scenarios:
//
// 1. Valid Path and Bucket Name.
// 2. Empty Database Path.
// 3. Empty Bucket Name.
func TestValidateTransactionRepositoryBoltDB(t *testing.T) {
	tests := []struct {
		name          string
		pathToDB      string
		bucketName    string
		expectedError error
	}{
		{
			name:          "Valid Path And Bucket Name",
			pathToDB:      "testdata/transaction_test.db",
			bucketName:    "transactions",
			expectedError: nil,
		},
		{
			name:          "Empty Database Path",
			pathToDB:      "",
			bucketName:    "transactions",
			expectedError: repository.ErrPathToDBAndBucketNameIsMandatory,
		},
		{
			name:          "Empty Bucket Name",
			pathToDB:      "testdata/transaction_test.db",
			bucketName:    "",
			expectedError: repository.ErrPathToDBAndBucketNameIsMandatory,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := repository.ValidateTransactionRepositoryBoltDB(tt.pathToDB, tt.bucketName)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
