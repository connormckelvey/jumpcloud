package main

import (
	"net/http"
)

type App struct {
	*http.ServeMux
	middleware []Middleware
	metrics    interface{}
}

func NewApp() *App {
	return &App{
		ServeMux:   http.NewServeMux(),
		middleware: []Middleware{},
	}
}

func (a *App) Handler() http.Handler {
	last := a.ServeMux.ServeHTTP
	for i := len(a.middleware) - 1; i >= 0; i-- {
		last = a.middleware[i](last)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		last(w, r)
	})
}

func (a *App) Use(mw ...Middleware) {
	a.middleware = append(a.middleware, mw...)
}

func (a *App) All(path string, handler http.HandlerFunc) {
	a.HandleFunc(path, handler)
}

func (a *App) Get(path string, handler http.HandlerFunc) {
	a.HandleFunc(path, withMethod(http.MethodGet)(handler))
}

func (a *App) Post(path string, handler http.HandlerFunc) {
	a.HandleFunc(path, withMethod(http.MethodPost)(handler))
}

func withMethod(method string) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != method {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}
			next(w, r)
		}
	}
}
