package pwhash

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"
)

func newTestLogger() *log.Logger {
	output := ioutil.Discard
	if testing.Verbose() {
		output = os.Stdout
	}
	return log.New(output, "", log.LstdFlags)
}

func TestPasswordHashing(t *testing.T) {
	app := NewTestApp(t, &Config{
		Logger: newTestLogger(),
	})
	defer app.Quit()

	tests := []struct {
		in       string
		expected string
	}{
		{"", "z4PhNX7vuL3xVChQ1m2AB9Yg5AULVxXcg/SpIdNs6c5H0NE8XYXysP+DGNKHfuwvY7kxvUdBeoGlODJ6+SfaPg=="},
		{"angryMonkey", "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="},
		{"superSecret", "pQaPqt7aC/CThmNsO8xnV+nkLfyJoyqpFzGzmvLivIpjmQXnvJqIULCUOpE+H1f3+p9laadfIkvAxMYZTAxnyQ=="},
	}

	for _, test := range tests {
		form := url.Values{"password": []string{test.in}}
		res, err := app.Client().PostForm("/hash", form)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != 200 {
			t.Errorf("Expected: %d, got: %d", 200, res.StatusCode)
		}
		body, _ := ioutil.ReadAll(res.Body)
		if string(body) != test.expected {
			t.Errorf("Expected: %s, got: %s", test.expected, body)
		}
	}
}

func TestResponseDelay(t *testing.T) {
	app := NewTestApp(t, &Config{
		Logger: newTestLogger(),
	})
	defer app.Quit()

	if app.config.listenAddr() != ":8080" {
		t.Errorf("Expected: %s, got: %s", ":8080", app.config.listenAddr())
	}

	tests := []struct {
		inMethod             string
		inPath               string
		inParams             url.Values
		expectedCode         int
		expectedResponseTime time.Duration
	}{
		{"GET", "/", nil, 404, 0},
		{"POST", "/hash", url.Values{"password": []string{"angryMonkey"}}, 200, 5},
		{"GET", "/stats", nil, 200, 0},
		{"GET", "/shutdown", nil, 202, 0},
	}

	for _, test := range tests {
		startTime := time.Now()
		var res *http.Response
		var err error
		switch test.inMethod {
		case "GET":
			res, err = app.Client().Get(test.inPath)
		case "POST":
			res, err = app.Client().PostForm(test.inPath, test.inParams)
		}

		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != test.expectedCode {
			t.Errorf("Expected: %d, got: %d", test.expectedCode, res.StatusCode)
		}

		actualResponseTime := time.Now().Sub(startTime) / time.Second
		if actualResponseTime != test.expectedResponseTime {
			t.Errorf("Expected: %d, got: %d", test.expectedResponseTime, actualResponseTime)
		}
	}
}

func TestGracefulShutdown(t *testing.T) {
	app := NewTestApp(t, &Config{
		Logger: newTestLogger(),
	})
	defer app.Quit()

	numRequests := 25
	lastRequestSent := make(chan bool)
	wg := &sync.WaitGroup{}

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			form := url.Values{"password": []string{"angryMonkey"}}
			res, err := app.Client().PostForm("/hash", form)
			if err != nil {
				t.Fatal(err)
			}

			if res.StatusCode != 200 {
				t.Errorf("Expected: %d, got: %d", 200, res.StatusCode)
			}

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if len(body) != 88 {
				t.Errorf("Expected: %d, got: %d", 88, len(body))
			}
		}(i)

		if i == numRequests-1 {
			close(lastRequestSent)
		}
	}

	<-lastRequestSent

	time.AfterFunc(1*time.Second, func() {
		_, err := app.Client().Get("/shutdown")
		if err != nil {
			t.Fatal(err)
		}
	})

	wg.Wait()
}

func TestStats(t *testing.T) {
	app := NewTestApp(t, &Config{
		Logger: newTestLogger(),
	})
	defer app.Quit()

	numRequests := 25
	wg := &sync.WaitGroup{}

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			form := url.Values{"password": []string{"angryMonkey"}}
			_, err := app.Client().PostForm("/hash", form)
			if err != nil {
				t.Fatal(err)
			}
		}(i)
	}
	wg.Wait()

	res, err := app.Client().Get("/stats")
	if err != nil {
		t.Fatal(err)
	}

	body, _ := ioutil.ReadAll(res.Body)
	stats := struct {
		Total   int `json:"total"`
		Average int `json:"average"`
	}{}
	json.Unmarshal(body, &stats)

	if stats.Total != numRequests {
		t.Errorf("Expected: %d, got: %d", numRequests, stats.Total)
	}
	if stats.Average <= 0 {
		t.Errorf("Expected: > 0 , got: %d", stats.Average)
	}
}
