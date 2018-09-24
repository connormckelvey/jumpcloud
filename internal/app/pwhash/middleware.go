package pwhash

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/connormckelvey/jumpcloud/pkg/httputil/middleware"
	"github.com/connormckelvey/jumpcloud/pkg/metrics"
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

func withDelay(amount time.Duration) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			delayedWriter := &delayedResponseWriter{
				ResponseWriter: w,
				timer:          time.NewTimer(amount),
				once:           sync.Once{},
			}
			next.ServeHTTP(delayedWriter, r)
		})
	}
}

type responseRecorderWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseRecorderWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func withLogging(logger *log.Logger) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &responseRecorderWriter{
				ResponseWriter: w,
			}
			next.ServeHTTP(rw, r)
			logger.Println(r.Method, r.URL.Path, rw.status, r.RemoteAddr, r.UserAgent())
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

type metricRecorderWriter struct {
	http.ResponseWriter
	*metrics.Collector
	startTime time.Time
	once      sync.Once
}

func (w *metricRecorderWriter) Write(p []byte) (n int, err error) {
	w.once.Do(func() {
		elapsedTime := time.Now().Sub(w.startTime)
		w.Collector.Observe(int64(elapsedTime / time.Microsecond))
	})
	return w.ResponseWriter.Write(p)
}

func withMetrics(name string) middleware.Middleware {
	collector := metrics.NewCollector(name)
	metrics.Register(collector)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &metricRecorderWriter{
				ResponseWriter: w,
				Collector:      collector,
				startTime:      time.Now(),
				once:           sync.Once{},
			}
			next.ServeHTTP(rw, r)
		})
	}
}
