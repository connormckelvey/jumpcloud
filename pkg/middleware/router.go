package middleware

import (
	"net/http"
)

type Router struct {
	routes map[string]map[string]Middleware
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]map[string]Middleware),
	}
}

func (rt *Router) register(path string, method string, mw Middleware) {
	if methods := rt.routes[path]; methods == nil {
		rt.routes[path] = make(map[string]Middleware)
	}
	rt.routes[path][method] = mw
}

// Maybe get/post methods should just be builders for a middleware that handles that shit
func (rt *Router) Get(path string, mw Middleware) *Router {
	rt.register(path, http.MethodGet, mw)
	return rt
}

func (rt *Router) Post(path string, mw Middleware) *Router {
	rt.register(path, http.MethodPost, mw)
	return rt
}

// Does't allow for plugable middleware... metrics exporter router for example..
// Currently would need to be added by the router... not the worst but..
func (rt *Router) AsMiddleware() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			path, method := r.URL.Path, r.Method
			var (
				methods    map[string]Middleware
				middleware Middleware
				ok         bool
			)

			if methods, ok = rt.routes[path]; !ok {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}

			if middleware, ok = methods[method]; !ok {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}

			middleware(next)(w, r)
		}
	}
}
