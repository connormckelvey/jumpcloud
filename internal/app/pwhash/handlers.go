package pwhash

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/connormckelvey/jumpcloud/pkg/hashutil/sha512"
)

func (a *Application) handleHash() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pr, pw := io.Pipe()

		go func() {
			defer pw.Close()

			hasher := sha512.NewStringWriter()
			hasher.WriteString(r.PostFormValue("password"))

			enc := base64.NewEncoder(base64.StdEncoding, pw)
			defer enc.Close()

			enc.Write(hasher.Sum(nil))
		}()

		w.WriteHeader(http.StatusOK)
		io.Copy(w, pr)
	})
}

func (a *Application) handleShutdown() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.Quit()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, http.StatusText(http.StatusOK))
	})
}
