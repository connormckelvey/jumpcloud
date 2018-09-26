package metrics

import (
	"testing"
)

func TestStore(t *testing.T) {
	tests := []struct {
		inName        string
		inDelta       int64
		expectedValue int64
	}{
		{"foo", 100, 100},
		{"bar", 10, 10},
		{"baz", 123, 123},
	}

	store := NewStore()

	for _, test := range tests {
		store.Register(NewInt64Collector(test.inName))
		collector, ok := store.Get(test.inName).(*Int64Collector)
		if !ok {
			t.Fatal()
		}
		collector.Add(test.inDelta)
		if collector.Value() != test.expectedValue {
			t.Errorf("Expected %d, go %d", test.expectedValue, collector.Value())
		}
	}
}
