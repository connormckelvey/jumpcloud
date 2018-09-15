package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Middleware provides a convenient way of composing handlers for HTTP
// requests. It returns a new http.HandlerFunc and should call `next` to
// invoke the next handler in the middleware pipeline.
type Middleware func(next http.HandlerFunc) http.HandlerFunc

func RequestLogging() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent()) // defer this
			next.ServeHTTP(w, r)
		}
	}
}

func WithTiming() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			next.ServeHTTP(w, r)
			totalTime := time.Now().Sub(startTime)
			fmt.Println("Total time:", totalTime)
		}
	}
}
