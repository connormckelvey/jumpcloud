package metrics

import (
	"sync/atomic"
)

type Collector struct {
	Name  string
	sum   *Metric
	count *Metric
}

func NewCollector(name string) *Collector {
	return &Collector{
		Name:  name,
		sum:   NewMetric(name + "_sum"),
		count: NewMetric(name + "_count"),
	}
}

func (c *Collector) Observe(value int64) {
	c.sum.Add(value)
	c.count.Add(1)
}

func (c *Collector) Count() int64 {
	return atomic.LoadInt64(c.count.Value)
}

func (c *Collector) Sum() int64 {
	return atomic.LoadInt64(c.sum.Value)
}

func (c *Collector) Average() int64 {
	count := c.Count()
	if count > 0 {
		return c.Sum() / count
	}
	return 0
}
