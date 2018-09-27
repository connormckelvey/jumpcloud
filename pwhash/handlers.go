package pwhash

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/connormckelvey/jumpcloud/middleware"
	"github.com/connormckelvey/jumpcloud/routing"
)

// Handler initializes the http.ServeMux router with handlers and middleware
// and returns a http.Handler interface for use by http.Server
func (a *Application) Handler() http.Handler {
	a.router.Handle("/hash", &routing.MethodHandler{
		Post: middleware.New(
			a.withShutdown(),
			a.withFormValidation(hashPasswordFormKey),
			a.withDelay(a.config.hashDelaySeconds()),
			a.withDurationMetrics(hashTimeMetricKey, time.Microsecond),
		).WrapFunc(a.handleHash),
	})
	a.router.Handle("/shutdown", &routing.MethodHandler{
		Get: middleware.New(
			a.withShutdown(),
		).WrapFunc(a.handleShutdown),
	})
	a.router.Handle("/stats", &routing.MethodHandler{
		GetFunc: a.handleStats,
	})
	return middleware.New(a.withTracing(), a.withLogging()).Wrap(a.router)
}

func (a *Application) handleHash(w http.ResponseWriter, r *http.Request) {
	a.waitGroup.Add(1)
	defer a.waitGroup.Done()

	hasher := sha512.New()
	_, err := io.WriteString(hasher, r.PostFormValue(hashPasswordFormKey))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	pr, pw := io.Pipe()
	go func() {
		enc := base64.NewEncoder(base64.StdEncoding, pw)
		defer pw.Close()
		defer enc.Close()
		enc.Write(hasher.Sum(nil))
	}()
	io.Copy(w, pr)
}

func (a *Application) handleShutdown(w http.ResponseWriter, r *http.Request) {
	a.Quit()
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, http.StatusText(http.StatusAccepted))
}

func (a *Application) handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		json.NewEncoder(w).Encode(a.metrics.Get(hashTimeMetricKey))
	}()
	io.Copy(w, pr)
}
