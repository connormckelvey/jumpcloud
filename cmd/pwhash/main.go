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

type reqCtxKey int

const (
	reqTimestamp reqCtxKey = iota
)

var (
	bindHost string
	bindPort int
)

func main() {
	flag.StringVar(&bindHost, "h", "", "bind server host")
	flag.IntVar(&bindPort, "p", 8080, "bind server port")
	flag.Parse()

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	router := http.NewServeMux()

	router.Handle("/hash", &rt.PathHandler{
		Post: &mw.Handler{
			Handler: hashPassword(),
			Middleware: []mw.Middleware{
				delayResponseWriter(5 * time.Second),
			},
		},
	})

	handler := &mw.Handler{
		Handler: router,
		Middleware: []mw.Middleware{
			requestTimestamp(),
			logging(logger),
		},
	}

	address := fmt.Sprintf("%s:%d", bindHost, bindPort)
	server := &http.Server{
		Addr:         address,
		Handler:      handler,
		ErrorLog:     logger,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	logger.Println("Server is listening at", address)
	if err := server.ListenAndServe(); err != nil {
		logger.Fatalf("Could not listen on %s: %v\n", address, err)
	}
}
