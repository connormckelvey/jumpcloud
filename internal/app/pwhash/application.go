package pwhash

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/connormckelvey/jumpcloud/pkg/httputil/middleware"
	"github.com/connormckelvey/jumpcloud/pkg/httputil/routing"
)

type Application struct {
	Logger   *log.Logger
	quitOnce *sync.Once
	quit     chan bool
}

func New(logger *log.Logger) *Application {
	return &Application{
		Logger:   logger,
		quitOnce: &sync.Once{},
		quit:     make(chan bool, 1),
	}
}

func (a *Application) Handler() http.Handler {
	router := http.NewServeMux()

	hashMiddleware := middleware.New(withFormValidation("password"),
		withDelay(5*time.Second))

	router.Handle("/hash", &routing.PathHandler{
		Post: hashMiddleware.Wrap(a.handleHash()),
	})

	router.Handle("/shutdown", &routing.PathHandler{
		Get: a.handleShutdown(),
	})

	return middleware.New(withLogging(a.Logger)).Wrap(router)
}

func (a *Application) Quit() {
	a.quitOnce.Do(func() {
		close(a.quit)
	})
}

func (a *Application) Wait() {
	<-a.quit
}
