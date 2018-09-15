package router

import "net/http"

type Router struct {
	*http.ServeMux
}

func (r *Router) All(path string, handler http.HandlerFunc) {
	r.HandleFunc(path, handler)
}

func (r *Router) Get(path string, handler http.HandlerFunc) {
	r.HandleFunc(path, withMethod(http.MethodGet)(handler))
}

func (r *Router) Post(path string, handler http.HandlerFunc) {
	r.HandleFunc(path, withMethod(http.MethodPost)(handler))
}
