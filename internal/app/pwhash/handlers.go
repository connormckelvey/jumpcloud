package pwhash

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/connormckelvey/jumpcloud/pkg/hashutil/sha512"
	"github.com/connormckelvey/jumpcloud/pkg/httputil/middleware"
	"github.com/connormckelvey/jumpcloud/pkg/httputil/routing"
)

var handleHash = &routing.PathHandler{
	Post: middleware.New(
		withFormValidation("password"),
		withDelay(5*time.Second),
	).Wrap(instance.handleHash()),
}

var handleShutdown = &routing.PathHandler{
	Get: instance.handleShutdown(),
}

func init() {
	router.Handle("/hash", handleHash)
	router.Handle("/shutdown", handleShutdown)
}

func (a *application) handleHash() http.Handler {
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

func (a *application) handleShutdown() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.quit()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, http.StatusText(http.StatusOK))
	})
}
