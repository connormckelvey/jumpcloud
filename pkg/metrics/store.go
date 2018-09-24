package metrics

import (
	"sync"
)

type Store struct {
	collectors map[string]*Collector
	mux        sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		collectors: make(map[string]*Collector),
		mux:        sync.RWMutex{},
	}
}

func (s *Store) Register(collector *Collector) {
	s.collectors[collector.Name] = collector
}

func (s *Store) Collector(name string) (*Collector, bool) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	c, exists := s.collectors[name]
	return c, exists
}

var defaultStore = NewStore()

func Register(collector *Collector) {
	defaultStore.Register(collector)
}

func FindCollector(name string) (*Collector, bool) {
	return defaultStore.Collector(name)
}
