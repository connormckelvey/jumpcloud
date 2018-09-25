package pwhash

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var okHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, http.StatusText(http.StatusOK))
})

func TestWithDelay(t *testing.T) {
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
		handler := withDelay(test.inAmount)(okHandler)

		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		startTime := time.Now()
		handler.ServeHTTP(rr, req)
		endTime := time.Now()

		actual := int(endTime.Sub(startTime) / time.Second)
		if actual != test.expected {
			t.Errorf("Expected: %d, got: %d", test.expected, actual)
		}
	}
}

func statusHandler(status int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		fmt.Fprintf(w, "%s", http.StatusText(status))
	})
}

func TestWithLogging(t *testing.T) {
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
		handler := withLogging(logger)(statusHandler(test.inStatus))

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

func TestWithFormValidation(t *testing.T) {
	tests := []struct {
		inKey        string
		inValue      string
		expectedCode int
		expectedBody string
	}{
		{"password", "angryMonkey", 200, "OK"},
		{"password", "", 200, "OK"},
		{"passphrase", "angryMonkey", 422, "Missing param: password"},
		{"", "", 422, "Missing param: password"},
	}

	for _, test := range tests {
		rr := httptest.NewRecorder()
		handler := withFormValidation("password")(okHandler)

		data := strings.NewReader(fmt.Sprintf("%s=%s", test.inKey, test.inValue))
		req, err := http.NewRequest(http.MethodPost, "/", data)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		handler.ServeHTTP(rr, req)

		if rr.Code != test.expectedCode {
			t.Errorf("Expected: %d, got: %d", test.expectedCode, rr.Code)
		}

		actualBody := strings.TrimSpace(rr.Body.String())
		if actualBody != test.expectedBody {

			t.Errorf("Expected: %v, got: %v", len(test.expectedBody), len(actualBody))
		}
	}
}
