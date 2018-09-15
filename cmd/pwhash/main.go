package main

import (
	"flag"
	"fmt"
	"net/http"
)

// Other stuff
// w.Header().Set("Content-Type", "text/plain; charset=utf-8")
// 		w.Header().Set("X-Content-Type-Options", "nosniff")
// GOOD EXAMPLE: https://gist.github.com/enricofoltran/10b4a980cd07cb02836f70a4ab3e72d7
// https://medium.com/statuscode/how-i-write-go-http-services-after-seven-years-37c208122831

var serverPort int

func init() {
	flag.IntVar(&serverPort, "p", 8080, "specify a port to use")
	flag.Parse()
}

func main() {
	app := NewApp()
	app.Use(RequestLogging(), WithTiming())

	app.Post("/hash", app.handleHashPassword())
	app.Get("/stats", app.handleMetrics())
	app.All("/", app.handleIndex())

	address := fmt.Sprintf(":%d", serverPort)
	server := &http.Server{Addr: address, Handler: app.Handler()}
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}
