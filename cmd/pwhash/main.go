package main

import (
	"flag"

	"github.com/connormckelvey/jumpcloud/pwhash"
)

var config pwhash.Config

func init() {
	flag.IntVar(&config.ListenPort, "p", 8080, "Port to listen on")
	flag.IntVar(&config.HashDelay, "d", 5, "Amount to delay hash response in seconds")
}

func main() {
	flag.Parse()
	pwhash.NewApplication(&config).Start()
}
