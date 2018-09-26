package pwhash

import (
	"encoding/json"
	"testing"
	"time"
)

func TestAverageDuration(t *testing.T) {
	tests := []struct {
		inName       string
		inUnits      time.Duration
		inCount      int
		inDelta      time.Duration
		expectedJSON string
	}{
		{"foo", time.Microsecond, 10, 100 * time.Microsecond, `{"total":10,"average":100}`},
		{"bar", time.Microsecond, 0, 100 * time.Microsecond, `{"total":0,"average":0}`},
		{"baz", time.Microsecond, 1, 100 * time.Second, `{"total":1,"average":100000000}`},
		{"qux", time.Second, 10, 100 * time.Second, `{"total":10,"average":100}`},
	}

	for _, test := range tests {
		collector := NewAverageDuration(test.inName, test.inUnits)
		for i := 0; i < test.inCount; i++ {
			collector.Observe(test.inDelta)
		}
		actualJSON, err := json.Marshal(collector)
		if err != nil {
			t.Fatal(err)
		}
		if string(actualJSON) != test.expectedJSON {
			t.Errorf("Expected: %s, got: %s", test.expectedJSON, actualJSON)
		}
	}
}
