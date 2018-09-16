package routing

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func methodHandler(method string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed),
				http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprintf(w, "%s", http.StatusText(http.StatusOK))
	})
}

func TestPathHandler(t *testing.T) {
	handler := &PathHandler{
		Get:  methodHandler(http.MethodGet),
		Post: methodHandler(http.MethodPost),
	}

	tests := []struct {
		method string
		status int
	}{
		{method: http.MethodGet, status: http.StatusOK},
		{method: http.MethodPost, status: http.StatusOK},
		{method: http.MethodPut, status: http.StatusMethodNotAllowed},
		{method: http.MethodPatch, status: http.StatusMethodNotAllowed},
		{method: http.MethodDelete, status: http.StatusMethodNotAllowed},
	}

	for _, test := range tests {
		req, err := http.NewRequest(test.method, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		actualStatus := rr.Code
		if actualStatus != test.status {
			t.Errorf("Expected: %d, got: %d \n", test.status, actualStatus)
		}

		actualBody := strings.TrimSpace(rr.Body.String())
		expectedBody := http.StatusText(test.status)
		if actualBody != expectedBody {
			t.Errorf("Expected: %s, got: %s \n", expectedBody, actualBody)
		}
	}
}
