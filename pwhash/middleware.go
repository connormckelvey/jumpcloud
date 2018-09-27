package pwhash

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/connormckelvey/jumpcloud/middleware"
)

type delayedResponseWriter struct {
	http.ResponseWriter
	timer *time.Timer
	once  sync.Once
}

func (d *delayedResponseWriter) Write(p []byte) (n int, err error) {
	d.once.Do(func() {
		<-d.timer.C
	})
	return d.ResponseWriter.Write(p)
}

// withDelay returns a middleware.Middleware for use by any http.Handler
// Using delayedResponseWriter.Write it waits once for the time.Duration `amount`
// before proxying Write calls to the embedded http.ResponseWriter.
func (a *Application) withDelay(amount time.Duration) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &delayedResponseWriter{
				ResponseWriter: w,
				timer:          time.NewTimer(amount),
				once:           sync.Once{},
			}
			next.ServeHTTP(rw, r)
		})
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// withLogging returns a middleware.Middleware for use by any http.Handler.
// It calls the `next` handler to capture http.StatusCodes via
// loggingResponseWriter.Write and finally logs request information using
// Application.logger
func (a *Application) withLogging() middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &loggingResponseWriter{
				ResponseWriter: w,
			}
			next.ServeHTTP(rw, r)
			traceID, ok := r.Context().Value(traceIDKey).(string)
			if !ok {
				traceID = "No-Trace-ID"
			}
			a.logger.Println(traceID, r.Method, r.URL.Path, rw.status, r.RemoteAddr, r.UserAgent())
		})
	}
}

// withTracing returns a middleware.Middleware for use by any http.Handler.
// It checks for a `X-Request-ID` used for tracing. If `X-Request-ID` is empty
// a new TraceID is created and add to the request's Context and as a Header
// on the response.
func (a *Application) withTracing() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			traceID := r.Header.Get("X-Request-ID")
			if traceID == "" {
				traceID = fmt.Sprintf("%d", time.Now().UnixNano())
			}
			ctx := context.WithValue(r.Context(), traceIDKey, traceID)
			w.Header().Set("X-Request-Id", traceID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// withFormValidation returns a middleware.Middleware for use by any http.Handler.
// It parses the request's PostForm and validates the list of provided `requiredParams`.
// If any `requiredParam`s are missing, the request is halted and an error is sent
// to the client.
func (a *Application) withFormValidation(requiredParams ...string) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			for _, param := range requiredParams {
				if _, ok := r.PostForm[param]; !ok {
					http.Error(w, fmt.Sprintf("Missing param: %s", param),
						http.StatusUnprocessableEntity)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

type metricRecorderWriter struct {
	http.ResponseWriter
	collector *AverageDuration
	startTime time.Time
	once      sync.Once
}

func (w *metricRecorderWriter) Write(p []byte) (n int, err error) {
	w.once.Do(func() {
		w.collector.Observe(time.Now().Sub(w.startTime))
	})
	return w.ResponseWriter.Write(p)
}

// withDurationMetrics returns a middleware.Middleware for use by any http.Handler.
// It creates and registers a new AverageDuration metrics.Collector using the provided
// `name` and `unit` of time. It calls AverageDuration.Observe using the time elapsed
// between the intital request and the first call to metricRecorderWriter.Write.
func (a *Application) withDurationMetrics(name string, unit time.Duration) middleware.Middleware {
	collector := NewAverageDuration(hashTimeMetricKey, time.Microsecond)
	a.metrics.Register(collector)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &metricRecorderWriter{
				ResponseWriter: w,
				collector:      collector,
				startTime:      time.Now(),
				once:           sync.Once{},
			}
			next.ServeHTTP(rw, r)
		})
	}
}

// withDurationMetrics returns a middleware.Middleware for use by any http.Handler.
// If Application is shutting down it halts the request and returns an error
// to the client.
func (a *Application) withShutdown() middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if a.InShutdown() {
				http.Error(w, http.StatusText(http.StatusServiceUnavailable),
					http.StatusServiceUnavailable)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
