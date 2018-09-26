package pwhash

import (
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

func (a *Application) withLogging() middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &loggingResponseWriter{
				ResponseWriter: w,
			}
			next.ServeHTTP(rw, r)
			a.logger.Println(r.Method, r.URL.Path, rw.status, r.RemoteAddr, r.UserAgent())
		})
	}
}

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
	collector *averageDuration
	startTime time.Time
	once      sync.Once
}

func (w *metricRecorderWriter) Write(p []byte) (n int, err error) {
	w.once.Do(func() {
		w.collector.Observe(time.Now().Sub(w.startTime))
	})
	return w.ResponseWriter.Write(p)
}

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
