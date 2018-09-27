package pwhash

import (
	"fmt"
	"log"
	"os"
	"time"
)

type contextKey int

const (
	traceIDKey contextKey = iota
)

const (
	hashTimeMetricKey   = "hashTime"
	hashPasswordFormKey = "password"
	defaultHashDelay    = 5
	defaultListenPort   = 8080
)

// Config provides configuration fields that are used when setting up and
// starting Application.
type Config struct {
	Logger     *log.Logger
	ListenPort int
	HashDelay  int
}

func (c *Config) hashDelaySeconds() time.Duration {
	if c.HashDelay == 0 {
		c.HashDelay = defaultHashDelay
	}
	return time.Duration(c.HashDelay) * time.Second
}

func (c *Config) listenAddr() string {
	if c.ListenPort == 0 {
		c.ListenPort = defaultListenPort
	}
	return fmt.Sprintf(":%d", c.ListenPort)
}

func (c *Config) logger() *log.Logger {
	if c.Logger != nil {
		return c.Logger
	}
	return log.New(os.Stdout, "", log.LstdFlags)
}
