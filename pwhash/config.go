package pwhash

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Config struct {
	Logger     *log.Logger
	ListenPort int
	HashDelay  int
}

func (c *Config) hashDelaySeconds() time.Duration {
	if c.HashDelay == 0 {
		c.HashDelay = 5
	}
	return time.Duration(c.HashDelay) * time.Second
}

func (c *Config) listenAddr() string {
	if c.ListenPort == 0 {
		c.ListenPort = 8080
	}
	return fmt.Sprintf(":%d", c.ListenPort)
}

func (c *Config) logger() *log.Logger {
	if c.Logger != nil {
		return c.Logger
	}
	return log.New(os.Stdout, "", log.LstdFlags)
}
