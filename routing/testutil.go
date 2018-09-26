package routing

import (
	"fmt"
	"net/http"
)

func methodHandlerFunc(method string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed),
				http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprintf(w, "%s", http.StatusText(http.StatusOK))
	}
}

func methodHandler(method string) http.Handler {
	return http.HandlerFunc(methodHandlerFunc(method))
}
