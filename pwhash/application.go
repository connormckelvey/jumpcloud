package pwhash

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
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
}

func NewApplication(config *Config) *Application {
	return &Application{
		config:    config,
		logger:    config.logger(),
		router:    http.NewServeMux(),
		waitGroup: &sync.WaitGroup{},
	}
}

func (a *Application) Start() {
	server := &http.Server{
		Addr:     a.config.listenAddr(),
		Handler:  a.Handler(),
		ErrorLog: a.logger,
	}

	a.waitGroup.Add(1)
	go func() {
		a.logger.Printf("Server is listening at %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal(err)
		}
	}()

	a.waitGroup.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.logger.Printf("Server is shutting down...\n")
	server.SetKeepAlivesEnabled(false)

	if err := server.Shutdown(ctx); err != nil {
		a.logger.Fatalf("Could not gracefully shutdown the server: %s\n", err)
	}
	a.logger.Println("Server stopped gracefully")
}

func (a *Application) InShutdown() bool {
	return atomic.LoadUint64(&a.inShutdown) == 1
}

func (a *Application) Quit() error {
	if a.InShutdown() {
		return errors.New("Application is already shutting down")
	}
	atomic.StoreUint64(&a.inShutdown, 1)
	a.waitGroup.Done()
	return nil
}
