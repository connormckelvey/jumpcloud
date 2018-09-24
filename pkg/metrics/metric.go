package metrics

import "sync/atomic"

type Metric struct {
	Name  string
	Value *int64
}

func NewMetric(name string) *Metric {
	var value int64
	return &Metric{
		Name:  name,
		Value: &value,
	}
}

func (m *Metric) Add(v int64) {
	atomic.AddInt64(m.Value, v)
}
