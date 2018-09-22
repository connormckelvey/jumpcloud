package main

import (
	"context"
	"flag"
	"net/http"
	"time"

	"github.com/connormckelvey/jumpcloud/internal/app/pwhash"
)

var (
	port   = flag.String("p", ":8080", "Port to listen on")
	logger = pwhash.Logger()
	done   = make(chan bool)
)

var server = &http.Server{
	Handler:  pwhash.Handler(),
	ErrorLog: logger,
}

func main() {
	flag.Parse()
	server.Addr = *port

	go func() {
		pwhash.Wait()
		defer close(done)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		logger.Printf("Server is shutting down...\n")
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %s\n", err)
		}
	}()

	logger.Printf("Server is listening at %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Printf("Could not listen at %s: %s\n", server.Addr, err)
	}

	<-done
	logger.Println("Server stopped gracefully")
}
