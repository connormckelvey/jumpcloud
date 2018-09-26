package pwhash

import (
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		inConfig           *Config
		expectedListenAddr string
		expectedHashDelay  time.Duration
	}{
		{&Config{}, ":8080", 5 * time.Second},
	}

	for _, test := range tests {
		logger := test.inConfig.logger()
		if logger == nil {
			t.Errorf("Expected a logger, not nil")
		}
		listenAddr := test.inConfig.listenAddr()
		if listenAddr != test.expectedListenAddr {
			t.Errorf("Expected :%s, got: %s", test.expectedListenAddr, listenAddr)
		}
		hashDelay := test.inConfig.hashDelaySeconds()
		if hashDelay != test.expectedHashDelay {
			t.Errorf("Expected :%s, got: %s", test.expectedHashDelay, hashDelay)
		}
	}
}
