package metrics

import "sync/atomic"

// Collector is an interface that represents an individual metric
// Name() is called during registration with Store.
type Collector interface {
	Name() string
}

// Int64Collector is a simple implementation of Collector. It stores a
// value as an int64.
type Int64Collector struct {
	name  string
	value int64
}

// NewInt64Collector is the constructor for Int64Collector. It accepts a string,
// `name` and returns *Int64Collector
func NewInt64Collector(name string) *Int64Collector {
	return &Int64Collector{
		name:  name,
		value: 0,
	}
}

// Name makes Int64Collector compatible with the Collector interface. It returns
// the name provided during initialization.
func (c *Int64Collector) Name() string {
	return c.name
}

// Add accepts an int64, `delta` and atomically adds `delta` to the value stored
// in *Int64Collector
func (c *Int64Collector) Add(delta int64) {
	atomic.AddInt64(&c.value, delta)
}

// Value atomically loads and returns the value stored in *Int64Collector
func (c *Int64Collector) Value() int64 {
	return atomic.LoadInt64(&c.value)
}
