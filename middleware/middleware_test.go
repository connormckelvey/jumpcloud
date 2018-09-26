package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChain(t *testing.T) {
	tests := []struct {
		in       Chain
		expected string
	}{
		{New(stackMw(0), stackMw(1)), "[0 1 -1]"},
		{New(stackMw(0), stackMw(1), stackMw(2)), "[0 1 2 -1]"},
		{New(stackMw(3), stackMw(2), stackMw(1)), "[3 2 1 -1]"},
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		callStack = []int{}
		handler := test.in.Wrap(stackHandler)

		res := httptest.NewRecorder()
		handler.ServeHTTP(res, req)

		actual := res.Body.String()

		if actual != test.expected {
			t.Errorf("Expected: %s, got %s \n", test.expected, actual)
		}
	}
}

var callStack []int

func stackMw(n int) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callStack = append(callStack, n)
			next.ServeHTTP(w, r)
		})
	}
}

var stackHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	callStack = append(callStack, -1)
	fmt.Fprintf(w, "%v", callStack)
})
