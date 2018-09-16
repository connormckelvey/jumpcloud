package middleware

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func statusHandler(status int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		fmt.Fprintf(w, "%s", http.StatusText(status))
	})
}

func TestLogging(t *testing.T) {
	tests := []struct {
		inPath      string
		inMethod    string
		inStatus    int
		inUserAgent string
		expected    string
	}{
		{"/", "GET", http.StatusOK, "Test/1.0", "GET / 200  Test/1.0\n"},
		{"/foo", "POST", http.StatusNotFound, "Test/1.0", "POST /foo 404  Test/1.0\n"},
	}

	for _, test := range tests {
		output := bytes.NewBuffer([]byte{})
		logger := log.New(output, "", log.LstdFlags)

		rr := httptest.NewRecorder()
		handler := Logging(logger)(statusHandler(test.inStatus))

		req, err := http.NewRequest(test.inMethod, test.inPath, nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("User-Agent", test.inUserAgent)
		handler.ServeHTTP(rr, req)

		logLine := output.String()
		if err != nil {
			t.Fatal(err)
		}

		logParts := strings.Split(logLine, " ")
		actual := strings.Join(logParts[2:], " ")

		if actual != test.expected {
			t.Errorf("Expected: %s, got: %s", test.expected, actual)
		}
	}
}
