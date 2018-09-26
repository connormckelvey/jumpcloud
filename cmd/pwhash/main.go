package main

import (
	"flag"
	"os"

	"github.com/connormckelvey/jumpcloud/pwhash"
)

var config pwhash.Config

func init() {
	flag.IntVar(&config.ListenPort, "p", 8080, "Port to listen on")
	flag.IntVar(&config.HashDelay, "d", 5, "Amount to delay hash response in seconds")
}

func main() {
	flag.Parse()
	if err := pwhash.NewApplication(&config).Start(); err != nil {
		os.Exit(1)
	}
}
