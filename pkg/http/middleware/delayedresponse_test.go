package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var okHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", http.StatusText(http.StatusOK))
})

func TestDelayedResponse(t *testing.T) {
	tests := []struct {
		inAmount time.Duration
		expected int
	}{
		{1 * time.Second, 1},
		{2 * time.Second, 2},
		{5 * time.Second, 5},
	}

	for _, test := range tests {
		rr := httptest.NewRecorder()
		handler := DelayedResponseWriter(test.inAmount)(okHandler)

		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		startTime := time.Now()
		handler.ServeHTTP(rr, req)
		endTime := time.Now()

		actual := int(endTime.Sub(startTime) / 1000000000)
		if actual != test.expected {
			t.Errorf("Expected: %d, got: %d", test.expected, actual)
		}
	}
}
