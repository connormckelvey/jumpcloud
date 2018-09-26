package metrics

import (
	"sync"
	"testing"
)

func TestInt64Collector(t *testing.T) {
	tests := []struct {
		inName        string
		inCount       int
		inDelta       int64
		expectedValue int64
	}{
		{"foo", 100, 1, 100},
		{"bar", 10, 10, 100},
		{"baz", 123, 2, 246},
	}

	for _, test := range tests {
		collector := NewInt64Collector(test.inName)
		wg := sync.WaitGroup{}
		for i := 0; i < test.inCount; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				collector.Add(test.inDelta)
			}(i)
		}
		wg.Wait()
		if collector.Value() != test.expectedValue {
			t.Errorf("Expected %d, go %d", test.expectedValue, collector.Value())
		}
	}
}
