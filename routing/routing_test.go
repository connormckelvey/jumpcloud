package routing

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func methodHandlerFunc(method string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed),
				http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprintf(w, "%s", http.StatusText(http.StatusOK))
	}
}

func methodHandler(method string) http.Handler {
	return http.HandlerFunc(methodHandlerFunc(method))
}

func TestMethodHandler(t *testing.T) {
	handler := &MethodHandler{
		Get:  methodHandler(http.MethodGet),
		Post: methodHandler(http.MethodPost),
	}

	handlerFunc := &MethodHandler{
		GetFunc:  methodHandlerFunc(http.MethodGet),
		PostFunc: methodHandlerFunc(http.MethodPost),
	}

	tests := []struct {
		handler *MethodHandler
		method  string
		status  int
	}{
		{handler, http.MethodGet, http.StatusOK},
		{handler, http.MethodPost, http.StatusOK},
		{handler, http.MethodPut, http.StatusMethodNotAllowed},
		{handler, http.MethodPatch, http.StatusMethodNotAllowed},
		{handler, http.MethodDelete, http.StatusMethodNotAllowed},
		{handlerFunc, http.MethodGet, http.StatusOK},
		{handlerFunc, http.MethodPost, http.StatusOK},
		{handlerFunc, http.MethodPut, http.StatusMethodNotAllowed},
		{handlerFunc, http.MethodPatch, http.StatusMethodNotAllowed},
		{handlerFunc, http.MethodDelete, http.StatusMethodNotAllowed},
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
