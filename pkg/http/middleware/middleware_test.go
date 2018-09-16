package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func ctxMw(key string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), key, time.Now())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

var ctxTestHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	foo := r.Context().Value("foo").(time.Time)
	bar := r.Context().Value("bar").(time.Time)
	fmt.Fprintf(w, "%d", foo.Sub(bar))
})

func TestHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := &Handler{
		Handler:    ctxTestHandler,
		Middleware: []Middleware{ctxMw("foo"), ctxMw("bar")},
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status: %d to be %d \n", rr.Code, http.StatusOK)
	}

	timeDiff, err := strconv.Atoi(rr.Body.String())
	if err != nil {
		t.Fatal(err)
	}

	if timeDiff >= 0 {
		t.Errorf("Expected timeDiff: %d to be < 0 \n", timeDiff)
	}
}
