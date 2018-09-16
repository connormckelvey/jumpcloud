package main

import (
	"fmt"
	"net/http"

	"github.com/connormckelvey/jumpcloud/pkg/stringhasher/sha512"
)

func hashPassword() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		password := r.PostForm.Get("password")
		if password == "" {
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity),
				http.StatusUnprocessableEntity)
			return
		}

		hashedPassword := sha512.HashAndBase64Encode(password)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintln(w, hashedPassword)
	})
}
