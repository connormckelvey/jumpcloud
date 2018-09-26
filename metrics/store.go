package metrics

import (
	"sync"
)

type Store struct {
	collectors map[string]Collector
	mux        sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		collectors: make(map[string]Collector),
		mux:        sync.RWMutex{},
	}
}

func (s *Store) Register(collector Collector) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.collectors[collector.Name()] = collector
}

func (s *Store) Get(name string) Collector {
	s.mux.RLock()
	defer s.mux.RUnlock()
	if c, exists := s.collectors[name]; exists {
		return c
	}
	return nil
}
