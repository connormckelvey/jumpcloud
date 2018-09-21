package middleware

import "net/http"

type Middleware func(next http.Handler) http.Handler

type Chain []Middleware

func New(middlewares ...Middleware) Chain {
	return Chain(middlewares)
}

func (c Chain) Wrap(handler http.Handler) http.Handler {
	for i := len(c) - 1; i >= 0; i-- {
		handler = c[i](handler)
	}
	return handler
}
