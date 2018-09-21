package pwhash

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/connormckelvey/jumpcloud/pkg/httputil/middleware"
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

func withDelay(amount time.Duration) middleware.Middleware {
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

type responseRecorderWriter struct {
	http.ResponseWriter
	Status int
}

func (w *responseRecorderWriter) WriteHeader(code int) {
	w.Status = code
	w.ResponseWriter.WriteHeader(code)
}

func withLogging(logger *log.Logger) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &responseRecorderWriter{
				ResponseWriter: w,
			}
			next.ServeHTTP(rw, r)
			logger.Println(r.Method, r.URL.Path, rw.Status, r.RemoteAddr, r.UserAgent())
		})
	}
}

func withFormValidation(requiredParams ...string) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, err.Error(), 422)
				return
			}
			for _, param := range requiredParams {
				if _, ok := r.PostForm[param]; !ok {
					http.Error(w, fmt.Sprintf("Missing param: %s", param), 422)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
