package routing

import (
	"net/http"
)

type PathHandler struct {
	http.Handler
	Post http.Handler
	Get  http.Handler
}

func (h *PathHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if http.MethodGet == r.Method && h.Get != nil {
		h.Get.ServeHTTP(w, r)
		return
	}

	if http.MethodPost == r.Method && h.Post != nil {
		h.Post.ServeHTTP(w, r)
		return
	}

	http.Error(w, http.StatusText(http.StatusMethodNotAllowed),
		http.StatusMethodNotAllowed)
}
