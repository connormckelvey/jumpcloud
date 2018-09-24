package pwhash

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/connormckelvey/jumpcloud/pkg/httputil/middleware"
)

type application struct {
	logger   *log.Logger
	handler  http.Handler
	quitOnce sync.Once
	quitChan chan bool
}

var (
	logger  = log.New(os.Stderr, "", log.LstdFlags)
	router  = http.DefaultServeMux
	handler = middleware.New(withLogging(logger)).Wrap(router)
)

var instance = &application{
	logger:   logger,
	handler:  handler,
	quitOnce: sync.Once{},
	quitChan: make(chan bool, 1),
}

func Logger() *log.Logger { return instance.logger }

func Handler() http.Handler { return instance.handler }

func Quit() {
	instance.quit()
}

func Wait() { instance.wait() }

func (a *application) quit() {
	a.quitOnce.Do(func() {
		close(a.quitChan)
	})
}

func (a *application) wait() {
	<-a.quitChan
}
