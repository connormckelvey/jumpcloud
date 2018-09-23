package metrics

import (
	"sync"
)

var defaultStore = NewStore()

type Store struct {
	data map[uint64]int64
	mux  sync.RWMutex
}

func NewStore() *Store {
	m := new(Store)
	m.Reset()
	return m
}

func (m *Store) Reset() {
	m.data = make(map[uint64]int64)
	m.mux = sync.RWMutex{}
}

func (m *Store) Put(key uint64, value int64) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.data[key] = value
}

func (m *Store) getUnsafe(key uint64) (int64, bool) {
	value, exists := m.data[key]
	if !exists {
		return 0, false
	}
	return value, true
}

func (m *Store) Get(key uint64) (int64, bool) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	return m.getUnsafe(key)
}

type UpdaterFunc func(int64) int64

func (m *Store) Update(key uint64, updater UpdaterFunc) {
	m.mux.Lock()
	defer m.mux.Unlock()

	current, _ := m.getUnsafe(key)
	m.data[key] = updater(current)
}
