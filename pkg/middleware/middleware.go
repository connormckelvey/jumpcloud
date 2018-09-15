package middleware

import (
	"fmt"
	"net/http"
)

// Middleware provides a convenient way of composing handlers for HTTP
// requests. It returns a new http.HandlerFunc and should call `next` to
// invoke the next handler in the middleware pipeline.
type Middleware func(next http.HandlerFunc) http.HandlerFunc

func composeMiddleware(mw ...Middleware) http.HandlerFunc {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Donezo!")
	}

	for i := len(mw) - 1; i >= 0; i-- {
		handlerFunc = mw[i](handlerFunc)
	}
	return handlerFunc
}

type Handler struct {
	handlerFunc http.HandlerFunc
}

func NewHandler(mw ...Middleware) *Handler {
	return &Handler{
		handlerFunc: composeMiddleware(mw...),
	}
}

func (mh Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mh.handlerFunc(w, r)
}
