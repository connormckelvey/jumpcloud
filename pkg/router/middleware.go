package router

import "net/http"

type handlerFuncMiddleware func(next http.HandlerFunc) http.HandlerFunc

func withMethod(method string) handlerFuncMiddleware {
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
