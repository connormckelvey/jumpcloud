package metrics

import "sync/atomic"

type Collector interface {
	Name() string
}

type Int64Collector struct {
	name  string
	value int64
}

func NewInt64Collector(name string) *Int64Collector {
	return &Int64Collector{
		name:  name,
		value: 0,
	}
}

func (c *Int64Collector) Name() string {
	return c.name
}

func (c *Int64Collector) Add(delta int64) {
	atomic.AddInt64(&c.value, delta)
}

func (c *Int64Collector) Value() int64 {
	return atomic.LoadInt64(&c.value)
}
