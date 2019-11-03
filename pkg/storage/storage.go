package storage

import (
	"8192bytes/util"
	"errors"
	"time"
)

var (
	// ErrorNotFound indicates that the requested data was not found
	ErrorNotFound = errors.New("data not found")
	// ErrorZeroLengthFileName indicates that the requested file name was zero length
	ErrorZeroLengthFileName = errors.New("zero length file name was provided")
	// ErrorZeroLengthData indicates that the data provided for storage was zero length
	ErrorZeroLengthData = errors.New("zero length data was provided")
)

const (
	// ExpireNever is the TTL value to use for an object that should never expire
	ExpireNever = time.Duration(-1)
)

// Storer wraps the store function
type Storer interface {
	Store(key string, data []byte) error
}

// Retriever wraps the retriever function
type Retriever interface {
	Retrieve(key string) ([]byte, error)
}

// StoreRetriever groups Storer and Retriever functions
type StoreRetriever interface {
	Storer
	Retriever
}

// Mediator handles sanitising data before storage and (if necessary) after retrieval
type Mediator struct {
	StoreRetriever StoreRetriever
}

// Store implements the storage mediator interface to StoreRetriever and allows sanitising data storage
func (sm *Mediator) Store(fileName string, data []byte) error {

	if len(fileName) == 0 {
		return ErrorZeroLengthFileName
	}

	if len(data) == 0 {
		return ErrorZeroLengthData
	}

	key, err := util.HashForPath(fileName)
	if err != nil {
		return err
	}

	return sm.StoreRetriever.Store(key, data)
}

// Load implments the storage mediate interface to StoreRetriever and allows sanitising data loads
func (sm *Mediator) Load(fileName string) ([]byte, error) {

	if len(fileName) == 0 {
		return nil, ErrorZeroLengthFileName
	}

	key, err := util.HashForPath(fileName)
	if err != nil {
		return nil, util.ErrorHashingFailed
	}

	data, err := sm.StoreRetriever.Retrieve(key)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, ErrorNotFound
	}

	return data, nil
}

// NewMediator instantiates a new Mediator for a gien StoreRetriever, allowing sanitisation of data before storage
func NewMediator(sr StoreRetriever) Mediator {
	return Mediator{
		StoreRetriever: sr,
	}
}
