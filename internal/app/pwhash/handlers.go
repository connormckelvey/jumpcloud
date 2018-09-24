package pwhash

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/connormckelvey/jumpcloud/pkg/httputil/middleware"
	"github.com/connormckelvey/jumpcloud/pkg/httputil/routing"
	"github.com/connormckelvey/jumpcloud/pkg/metrics"
)

const (
	hashTimeMetricKey = "hashTime"
)

var handleHash = &routing.PathHandler{
	Post: middleware.New(
		withFormValidation("password"),
		withDelay(5*time.Second),
		withMetrics(hashTimeMetricKey),
	).Wrap(instance.handleHash()),
}

var handleShutdown = &routing.PathHandler{
	Get: instance.handleShutdown(),
}

var handleStats = &routing.PathHandler{
	Get: instance.handleStats(),
}

func init() {
	router.Handle("/hash", handleHash)
	router.Handle("/shutdown", handleShutdown)
	router.Handle("/stats", handleStats)
}

func (a *application) handleHash() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pr, pw := io.Pipe()
		go func() {
			defer pw.Close()

			enc := base64.NewEncoder(base64.StdEncoding, pw)
			defer enc.Close()

			hasher := sha512.New()
			io.WriteString(hasher, r.PostFormValue("password"))

			enc.Write(hasher.Sum(nil))
		}()

		io.Copy(w, pr)
	})
}

func (a *application) handleShutdown() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.quit()
		fmt.Fprintln(w, http.StatusText(http.StatusOK))
	})
}

func (a *application) handleStats() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pr, pw := io.Pipe()
		go func() {
			defer pw.Close()

			collector, _ := metrics.FindCollector(hashTimeMetricKey)

			json.NewEncoder(w).Encode(map[string]int64{
				"total":   collector.Count(),
				"average": collector.Average(),
			})
		}()
		io.Copy(w, pr)
	})
}
