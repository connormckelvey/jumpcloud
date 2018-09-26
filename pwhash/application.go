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

const (
	hashTimeMetricKey   = "hashTime"
	hashPasswordFormKey = "password"
)

type Application struct {
	config     *Config
	logger     *log.Logger
	router     *http.ServeMux
	inShutdown uint64
	waitGroup  *sync.WaitGroup
	metrics    *metrics.Store
	errChan    chan error
}

func NewApplication(config *Config) *Application {
	return &Application{
		config:    config,
		logger:    config.logger(),
		router:    http.NewServeMux(),
		waitGroup: &sync.WaitGroup{},
		metrics:   metrics.NewStore(),
		errChan:   make(chan error),
	}
}

func (a *Application) Start() error {
	server := &http.Server{
		Addr:     a.config.listenAddr(),
		Handler:  a.Handler(),
		ErrorLog: a.logger,
	}

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)
	go func() {
		<-signals
		a.Quit()
	}()

	a.waitGroup.Add(1)
	go func() {
		a.logger.Printf("Server is listening at %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.QuitWithError(err)
		}
	}()

	a.waitGroup.Wait()
	if err := <-a.errChan; err != nil {
		a.logger.Println(err)
		return err
	}

	a.logger.Printf("Server shutting down...\n")
	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(context.Background()); err != nil {
		a.logger.Println(err)
		return err
	}

	finalMetrics, _ := json.Marshal(a.metrics.Get(hashTimeMetricKey))
	a.logger.Printf("Stats: %s \n", finalMetrics)

	a.logger.Println("Server stopped gracefully")
	return nil
}

func (a *Application) InShutdown() bool {
	return atomic.LoadUint64(&a.inShutdown) == 1
}

func (a *Application) Quit() {
	a.QuitWithError(nil)
}

func (a *Application) QuitWithError(err error) {
	if !a.InShutdown() {
		atomic.StoreUint64(&a.inShutdown, 1)
		a.waitGroup.Done()
		a.errChan <- err
	}
}
