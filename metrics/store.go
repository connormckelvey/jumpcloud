package metrics

import (
	"sync"
)

// Store holds references to Collectors in a Store.Store to provide a simple
// way of accessing a Collector by its name. It utilizes sync.RWMutex safe
// concurrent access and updates for the Store.collectors map.
type Store struct {
	collectors map[string]Collector
	mux        sync.RWMutex
}

// NewStore initializes and empty Store.collectors map and sync.RWMutex{}
// and returns a new *Store
func NewStore() *Store {
	return &Store{
		collectors: make(map[string]Collector),
		mux:        sync.RWMutex{},
	}
}

// Register accepts a Collector and stores it in the Store.collectors map.
// Registration typically happens from the main goroutine, but a Write lock
// is used while registering for added safety.
func (s *Store) Register(collector Collector) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.collectors[collector.Name()] = collector
}

// Get accepts a string, `name` and returns the associated Collector from
// Store.collectors. It uses a Read lock for safe concurrent access. If the
// Collector with the provided name is not found, nil is returned.
func (s *Store) Get(name string) Collector {
	s.mux.RLock()
	defer s.mux.RUnlock()
	if c, exists := s.collectors[name]; exists {
		return c
	}
	return nil
}
