package middleware

import (
	"log"
	"net/http"
)

type responseRecorderWriter struct {
	http.ResponseWriter
	Status int
}

func (w *responseRecorderWriter) WriteHeader(code int) {
	w.Status = code
	w.ResponseWriter.WriteHeader(code)
}

func Logging(logger *log.Logger) Middleware {
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
