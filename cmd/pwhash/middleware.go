package main

import (
	"context"
	"log"
	"net/http"
	"time"

	mw "github.com/connormckelvey/jumpcloud/pkg/http/middleware"
)

func logging(logger *log.Logger) mw.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				logger.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func requestTimestamp() mw.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), reqTimestamp, time.Now())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type delayedResponseWriter struct {
	shouldDelay bool
	*time.Timer
	http.ResponseWriter
}

func (d *delayedResponseWriter) Write(p []byte) (n int, err error) {
	if d.shouldDelay {
		<-d.Timer.C
		d.shouldDelay = false
	}
	return d.ResponseWriter.Write(p)
}

func delayResponseWriter(amount time.Duration) mw.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqReceivedTime, ok := r.Context().Value(reqTimestamp).(time.Time)
			if !ok {
				reqReceivedTime = time.Now()
			}

			delayDuration := reqReceivedTime.Add(amount).Sub(time.Now())
			delayedWriter := &delayedResponseWriter{
				shouldDelay:    true,
				Timer:          time.NewTimer(delayDuration),
				ResponseWriter: w,
			}
			next.ServeHTTP(delayedWriter, r)
		})
	}
}
