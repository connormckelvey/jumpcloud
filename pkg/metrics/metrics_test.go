package metrics

import (
	"sync"
	"testing"
)

func TestCounter(t *testing.T) {
	c := NewCounter()
	n := int64(10000)

	wg := sync.WaitGroup{}
	for i := int64(0); i < n; i++ {
		wg.Add(1)
		go func() {
			c.Inc()
			wg.Done()
		}()
	}
	wg.Wait()

	actualCount := c.Count()
	if actualCount != n {
		t.Errorf("Expected: %d, got: %d", n, actualCount)
	}
}

func TestSum(t *testing.T) {
	s := NewSum()
	n := int64(10000)

	wg := sync.WaitGroup{}
	for i := int64(0); i < n; i++ {
		wg.Add(1)
		go func(i int64) {
			s.Add(i)
			wg.Done()
		}(i)
	}
	wg.Wait()

	actualSum := s.Sum()
	m := n - 1
	expectedSum := ((m * m) + m) / 2
	if actualSum != expectedSum {
		t.Errorf("Expected: %d, got: %d", expectedSum, actualSum)
	}
}
