package metrics

import "sync/atomic"

type Sum interface {
	Add(v int64)
	Value() int64
}

type StandardSum struct {
	sum int64
}

func NewSum() Sum {
	return &StandardSum{0}
}

func (s *StandardSum) Add(v int64) {
	atomic.AddInt64(&s.sum, v)
}

func (s *StandardSum) Value() int64 {
	return atomic.LoadInt64(&s.sum)
}
