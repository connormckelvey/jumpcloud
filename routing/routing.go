package routing

import (
	"net/http"
)

// MethodHandler implements http.Handler and provides a declarative way
// to specify http.Handlers for HTTP Methods.
type MethodHandler struct {
	Post     http.Handler
	PostFunc func(w http.ResponseWriter, r *http.Request)
	Get      http.Handler
	GetFunc  func(w http.ResponseWriter, r *http.Request)
}

func (h *MethodHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if h.Get != nil {
			h.Get.ServeHTTP(w, r)
		} else if h.GetFunc != nil {
			http.HandlerFunc(h.GetFunc).ServeHTTP(w, r)
		}
		return
	case http.MethodPost:
		if h.Post != nil {
			h.Post.ServeHTTP(w, r)
		} else if h.PostFunc != nil {
			http.HandlerFunc(h.PostFunc).ServeHTTP(w, r)
		}
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed)
	}

}
