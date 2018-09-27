package middleware

import "net/http"

// Middleware is an interface that allows for composing http.Handlers
// into a middleware Chain. A Middleware is a function that accepts
// the `next` http.Handler in the Chain and returns and http.Handler.
type Middleware func(next http.Handler) http.Handler

// HandlerFunc is a Type Alias for a function that accepts an http.ResponseWrite
// and a *http.Request. Chain.WrapFunc accepts this type and converts it to an
// http.HandlerFunc.
type HandlerFunc = func(w http.ResponseWriter, r *http.Request)

// Chain is a type slice of Middleware. It contains Wrap and WrapFunc methods
// which perform the actual composition of the Middleware.
type Chain []Middleware

// New returns a new Chain from the slice of Middleware it receives.
func New(middlewares ...Middleware) Chain {
	return Chain(middlewares)
}

// Wrap loops backwards through the Chain of Middleware so that each Middleware is
// called in the order was passed in to New. It accepts `handler` which is used as
// the final http.Handler in the Middleware Chain
func (c Chain) Wrap(handler http.Handler) http.Handler {
	for i := len(c) - 1; i >= 0; i-- {
		handler = c[i](handler)
	}
	return handler
}

// WrapFunc accepts a HandlerFunc for a clean way of using functions instead of
// http.Handlers or http.HandlerFuncs. It converts the HandlerFunc into a
// http.HandlerFunc before calling Chain.Wrap
func (c Chain) WrapFunc(handlerFunc HandlerFunc) http.Handler {
	return c.Wrap(http.HandlerFunc(handlerFunc))
}
