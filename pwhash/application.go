package pwhash

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"

	"github.com/connormckelvey/jumpcloud/metrics"
)

// Application is the main container that stores dependencies and configuration
// information required by handlers, middledware, etc.
type Application struct {
	config     *Config
	logger     *log.Logger
	router     *http.ServeMux
	inShutdown uint64
	waitGroup  *sync.WaitGroup
	metrics    *metrics.Store
	errChan    chan error
}

// NewApplication returns an *Application built from the provided Config struct.
func NewApplication(config *Config) *Application {
	return &Application{
		config:    config,
		logger:    config.logger(),
		router:    http.NewServeMux(),
		waitGroup: &sync.WaitGroup{},
		metrics:   metrics.NewStore(),
		errChan:   make(chan error, 1),
	}
}

// Start begins the http.Server and waits for the server to stop. It returns an
// error if http.Server.ListenAndServe returns an error (such as port already
// in use), or if Application.QuitError is called with a non-nil error
func (a *Application) Start() error {
	server := &http.Server{
		Addr:     a.config.listenAddr(),
		Handler:  a.Handler(),
		ErrorLog: a.logger,
	}

	// Setup a signal chan to preform a graceful shutdown on os.Interupt
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)
	go func() {
		<-signals
		a.Quit()
	}()

	// Start the http.Server and add to Application.waitGroup to enable calling
	// Application.waitGroup.Add from non-main goroutines
	a.waitGroup.Add(1)
	go func() {
		a.logger.Printf("Server is listening at %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.QuitWithError(err)
		}
	}()

	// When Application.waitGroup is done, Application.Quit has been called and
	// all /hash requests have been completed
	a.waitGroup.Wait()

	// Wait for an error on Application.errChan. This allows Application.Start
	// to return an error from goroutines.
	if err := <-a.errChan; err != nil {
		a.logger.Println(err)
		return err
	}

	// Start the actual shutdown of the http.Server
	a.logger.Printf("Server shutting down...\n")
	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(context.Background()); err != nil {
		a.logger.Println(err)
		return err
	}

	// last but not least, log the final state of the hash stats
	finalMetrics, _ := json.Marshal(a.metrics.Get(hashTimeMetricKey))
	a.logger.Printf("Stats: %s \n", finalMetrics)

	a.logger.Println("Server stopped gracefully")
	return nil
}

// InShutdown returns whether or not Application is in the process of
// shutting down the http.Server
func (a *Application) InShutdown() bool {
	return atomic.LoadUint64(&a.inShutdown) == 1
}

// Quit calls Application.QuitWithError with a nil value indicating that the Application
// is quitting under expected circumstances. Called by the Application.handleShutdown
// http.Handler, and used in tearing down Application in tests.
func (a *Application) Quit() {
	a.QuitWithError(nil)
}

// QuitWithError starts the shutdown process. Application.waitGroup.Done is called to
// offset the Application.waitGroup.Add called in Application.Start. The error
// passed in is received in Application.Start.
func (a *Application) QuitWithError(err error) {
	if !a.InShutdown() {
		atomic.StoreUint64(&a.inShutdown, 1)
		a.waitGroup.Done()
		a.errChan <- err
	}
}
