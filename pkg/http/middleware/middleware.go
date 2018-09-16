package middleware

import "net/http"

type Middleware func(next http.Handler) http.Handler

func compose(mw ...Middleware) Middleware {
	return func(end http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var last = end
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}
			last.ServeHTTP(w, r)
		})
	}
}

type Handler struct {
	http.Handler
	Middleware []Middleware
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	compose(h.Middleware...)(h.Handler).ServeHTTP(w, r)
}
