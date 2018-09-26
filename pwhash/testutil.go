package pwhash

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func statusHandler(status int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		fmt.Fprintf(w, "%s", http.StatusText(status))
	})
}

var okHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, http.StatusText(http.StatusOK))
})

type TestApp struct {
	*Application
	URL     string
	errChan chan error
}

func NewTestApp(t *testing.T, config *Config) *TestApp {
	testApp := &TestApp{
		Application: NewApplication(config),
		errChan:     make(chan error),
	}
	listenAddr := testApp.config.listenAddr()
	testApp.URL = fmt.Sprintf("http://%s", listenAddr)
	go func() {
		defer close(testApp.errChan)
		if err := testApp.Start(); err != nil {
			testApp.errChan <- err
		}
	}()
	for {
		select {
		case err := <-testApp.errChan:
			if err != nil {
				testApp.Quit()
				t.Fatal(err)
			}
		default:
		}
		if _, err := testApp.Client().Get("/"); err == nil {
			break
		}
	}
	return testApp
}

type TestAppClient struct {
	http.Client
	*TestApp
}

func (a *TestApp) Client() *TestAppClient {
	return &TestAppClient{
		TestApp: a,
	}
}

func (c *TestAppClient) Get(path string) (*http.Response, error) {
	return c.Client.Get(c.TestApp.URL + path)
}

func (c *TestAppClient) PostForm(path string, data url.Values) (*http.Response, error) {
	return c.Client.PostForm(c.TestApp.URL+path, data)
}
