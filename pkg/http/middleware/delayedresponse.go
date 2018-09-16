package middleware

import (
	"net/http"
	"time"
)

type delayedResponseWriter struct {
	http.ResponseWriter
	*time.Timer
	shouldDelay bool
}

func (d *delayedResponseWriter) Write(p []byte) (n int, err error) {
	if d.shouldDelay {
		<-d.Timer.C
		d.shouldDelay = false
	}
	return d.ResponseWriter.Write(p)
}

func DelayedResponseWriter(amount time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			delayedWriter := &delayedResponseWriter{
				shouldDelay:    true,
				Timer:          time.NewTimer(amount),
				ResponseWriter: w,
			}
			next.ServeHTTP(delayedWriter, r)
		})
	}
}
