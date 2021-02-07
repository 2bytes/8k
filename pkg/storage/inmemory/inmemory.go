package inmemory

import (
	"errors"
	"fmt"
	"time"
)

const (
	// ExpireNever is the TTL value to use for an object that should never expire
	ExpireNever = time.Duration(-1)
)

var (
	// ErrorKeyUsed is returned for an attempt to store an object using an existing key is made
	ErrorKeyUsed = errors.New("the key specified is in use")

	// ErrorStorageFull is returned when the max item limit for object storage has been reached
	ErrorStorageFull = errors.New("storage limit has been reached")
)

// Object is the structure that holds an item that will be stored in memory. It can hold a TTL so that
type Object struct {
	Data    []byte
	Expires *time.Time
}

// Storage implements StoreRetriever for inmemory storage
type Storage struct {
	MaxItems  int
	TTL       time.Duration
	dataStore map[string]Object
}

// Store implements the Store function to satisfy StoreRetriever
func (s *Storage) Store(key string, data []byte) error {

	if len(s.dataStore) >= s.MaxItems {
		return ErrorStorageFull
	}

	if !s.Contains(key) {
		exp := time.Now().Add(s.TTL)
		s.dataStore[key] = Object{Data: data, Expires: &exp}
		fmt.Printf("Added value, size: %d. Status: %d/%d\n", len(data), len(s.dataStore), s.MaxItems)
		return nil
	}
	return ErrorKeyUsed
}

// Retrieve implements the Retrieve function to satisfy StoreRetriever
func (s *Storage) Retrieve(key string) ([]byte, error) {
	return s.dataStore[key].Data, nil
}

// Contains implements the Contains function to satisfy StoreRetriever
func (s *Storage) Contains(key string) bool {
	_, ok := s.dataStore[key]
	return ok
}

func (s *Storage) scrub() {
	for k, val := range s.dataStore {
		if val.Expires != nil && val.Expires.Before(time.Now()) {
			delete(s.dataStore, k)
			fmt.Printf("Deleted expired value. Status %d/%d\n", len(s.dataStore), s.MaxItems)
		}
	}
}

func (s *Storage) scrubber(interval time.Duration) {

	if interval > 0 {

		fmt.Printf("Starting scrubber with interval of %s\n", interval)

		ticker := time.NewTicker(interval)

		for range ticker.C {
			s.scrub()
		}
	}
}

// New initialises a new in memory store
// expiry sets the time before an object is deleted
// scrubInterval is the interval at which the scrubber will run to  clean up expired items
func New(expiry time.Duration, scrubInterval time.Duration, maxItems int) *Storage {

	fmt.Printf("Starting in-memory storage with TTL of %s\n", expiry)

	stg := &Storage{
		MaxItems:  maxItems,
		TTL:       time.Duration(expiry),
		dataStore: make(map[string]Object, maxItems),
	}

	if expiry > 0 && scrubInterval > 0 {
		go stg.scrubber(scrubInterval)
	}

	return stg
}
