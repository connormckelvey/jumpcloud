package middleware

import "net/http"

type Middleware func(next http.Handler) http.Handler
type HandlerFunc = func(w http.ResponseWriter, r *http.Request)

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

func (c Chain) WrapFunc(handlerFunc HandlerFunc) http.Handler {
	return c.Wrap(http.HandlerFunc(handlerFunc))
}
