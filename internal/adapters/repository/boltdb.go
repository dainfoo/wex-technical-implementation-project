package repository

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dainfoo/wex-technical-implementation-project/internal/core/domain"
	"github.com/google/uuid"
	"github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"go.etcd.io/bbolt"
)

// This file contains the implementation of the TransactionRepository interface using BoltDB.

// Activate the jsoniter library to decode the Treasury API response.
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// TransactionRepositoryBoltDB represents a BoltDB database with a bucket name to store transactions
// and a mutex to manage concurrent access to the database.
type TransactionRepositoryBoltDB struct {
	boltDB     *bbolt.DB
	bucketName string
	rwMutex    sync.RWMutex
}

// NewTransactionRepositoryBoltDB creates a new TransactionRepositoryBoltDB instance with input validation.
func NewTransactionRepositoryBoltDB(pathToDB string, bucketName string) (*TransactionRepositoryBoltDB, error) {
	pathToDB = strings.TrimSpace(pathToDB)
	bucketName = strings.TrimSpace(bucketName)

	if err := ValidateTransactionRepositoryBoltDB(pathToDB, bucketName); err != nil {
		return nil, err
	}

	// Ensures the directory exists or create it if it doesn't
	dir := filepath.Dir(pathToDB)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Error().Err(err).Msg("failed to create database directory")
		return nil, ErrDatabaseDirectoryCouldNotBeCreated
	}

	// Open the BoltDB database file with read-write permissions or create it if it doesn't exist
	boltDB, err := bbolt.Open(pathToDB, os.FileMode(0666), nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to create of open the database file")
		return nil, ErrCreateOpenDatabaseFile
	}

	// Ensures the bucket exists, or create it if it doesn't
	err = boltDB.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to create the bucket")
		return nil, ErrCreateBucket
	}

	return &TransactionRepositoryBoltDB{
		boltDB:     boltDB,
		bucketName: bucketName,
	}, nil
}

// SaveTransaction implements the SaveTransaction method of the TransactionRepository interface for BoltDB.
func (r *TransactionRepositoryBoltDB) SaveTransaction(transaction domain.Transaction) error {
	// Get a write lock to ensure exclusive access to the database
	// Only one transaction can be saved at a time to prevent deadlocks
	r.rwMutex.Lock()
	// Release the write lock after the function execution
	defer r.rwMutex.Unlock()

	return r.boltDB.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(r.bucketName))
		if bucket == nil {
			log.Error().
				Str("bucket", r.bucketName).
				Msg("bucket not found in BoltDB")
			return ErrBucketNotFound
		}

		transactionJsonData, err := json.Marshal(transaction)
		if err != nil {
			log.Error().
				Err(err).
				Str("transaction_id", transaction.ID.String()).
				Msg("failed to marshal transaction data")
			return err
		}

		err = bucket.Put([]byte(transaction.ID.String()), transactionJsonData)
		if err != nil {
			log.Error().
				Err(err).
				Str("transaction_id", transaction.ID.String()).
				Msg("failed to save the transaction")
		}
		return err
	})
}

// FindTransaction implements the FindTransaction method of the TransactionRepository interface for BoltDB.
func (r *TransactionRepositoryBoltDB) FindTransaction(id uuid.UUID) (*domain.Transaction, error) {
	// Get a read lock to ensure shared read access to the database
	r.rwMutex.RLock()
	// Release the read lock after the function execution
	defer r.rwMutex.RUnlock()

	var transaction domain.Transaction
	err := r.boltDB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(r.bucketName))
		if bucket == nil {
			log.Error().
				Str("bucket", r.bucketName).
				Msg("bucket not found in BoltDB")
			return ErrBucketNotFound
		}

		transactionJsonData := bucket.Get([]byte(id.String()))
		if transactionJsonData == nil {
			log.Warn().
				Str("transaction_id", id.String()).
				Msg("transaction not found in BoltDB")
			return ErrTransactionNotFound
		}

		err := json.Unmarshal(transactionJsonData, &transaction)
		if err != nil {
			log.Error().
				Err(err).
				Str("transaction_id", id.String()).
				Msg("failed to unmarshal transaction data")
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// GetBoltDB returns the BoltDB instance.
func (r *TransactionRepositoryBoltDB) GetBoltDB() *bbolt.DB {
	return r.boltDB
}

// Close closes the BoltDB database file.
func (r *TransactionRepositoryBoltDB) Close() error {
	return r.boltDB.Close()
}

// ValidateTransactionRepositoryBoltDB validates the path to the BoltDB database file and the bucket name
// for the TransactionRepositoryBoltDB struct.
func ValidateTransactionRepositoryBoltDB(pathToDB string, bucketName string) error {
	// Validate the database file path and bucket name emptiness: must not be empty
	if pathToDB == "" || bucketName == "" {
		return ErrPathToDBAndBucketNameIsMandatory
	}
	return nil
}
