package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/connormckelvey/jumpcloud/pkg/stringhasher/sha512"
)

func (a *App) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func (a *App) handleMetrics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "FUCK METRICS")
	}
}

func hashPassword(password string) string {
	hasher := sha512.New()
	hasher.WriteString(password)

	return hasher.String()
}

func (a *App) handleHashPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		password := r.PostForm.Get("password")
		hasher := sha512.New()
		hasher.WriteString(password)

		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s\n", hasher.String())
	}
}
