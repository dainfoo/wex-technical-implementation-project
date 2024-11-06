package repository

import "errors"

// This file defines error variables related to the BoltDB repository in the repository layer.

var (
	// ErrPathToDBAndBucketNameIsMandatory is returned when the database file path and/or the bucket name are empty.
	ErrPathToDBAndBucketNameIsMandatory = errors.New("the database file path and the bucket name are mandatory")

	// ErrDatabaseDirectoryCouldNotBeCreated is returned when the database file directory could not be created.
	ErrDatabaseDirectoryCouldNotBeCreated = errors.New("the database file directory could not be created")

	// ErrCreateOpenDatabaseFile is returned when the database file could not be created or opened.
	ErrCreateOpenDatabaseFile = errors.New("the database file could not be created or opened")

	// ErrCreateBucket is returned when the database bucket could not be created.
	ErrCreateBucket = errors.New("the database bucket could not be created")

	// ErrBucketNotFound is returned when the bucket is not found.
	ErrBucketNotFound = errors.New("bucket not found")

	// ErrTransactionNotFound is returned when the transaction is not found.
	ErrTransactionNotFound = errors.New("transaction not found")
)
