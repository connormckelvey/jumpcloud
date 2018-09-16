package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	mw "github.com/connormckelvey/jumpcloud/pkg/http/middleware"
	rt "github.com/connormckelvey/jumpcloud/pkg/http/routing"
)

func main() {
	var bindPort int
	flag.IntVar(&bindPort, "p", 8080, "bind server port")
	flag.Parse()

	address := fmt.Sprintf(":%d", bindPort)
	logger := log.New(os.Stdout, "", log.LstdFlags)
	router := http.NewServeMux()

	handler := &mw.Handler{
		Handler: router,
		Middleware: []mw.Middleware{
			mw.Logging(logger),
		},
	}

	server := &http.Server{
		Addr:         address,
		Handler:      handler,
		ErrorLog:     logger,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	router.Handle("/hash", &rt.PathHandler{
		Post: &mw.Handler{
			Handler: hashPassword(),
			Middleware: []mw.Middleware{
				mw.DelayedResponseWriter(5 * time.Second),
			},
		},
	})

	logger.Println("Server is listening at", address)
	if err := server.ListenAndServe(); err != nil {
		logger.Fatalf("Could not listen on %s: %v\n", address, err)
	}
}
