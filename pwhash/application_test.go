package pwhash

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"
)

func init() {
	*logger = *log.New(ioutil.Discard, "", log.LstdFlags)
}

func TestApplicationHandler(t *testing.T) {
	tests := []struct {
		inMethod             string
		inPath               string
		inParams             url.Values
		expectedCode         int
		expectedBody         string
		expectedResponseTime time.Duration
	}{
		{"GET", "/", nil, 404, "404 page not found", 0},
		{"GET", "/shutdown", nil, 200, "OK", 0},
		{http.MethodPost, "/hash", url.Values{"password": []string{"angryMonkey"}}, 200, "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==", 5},
	}

	server := httptest.NewServer(Handler())

	for _, test := range tests {
		client := server.Client()
		startTime := time.Now()

		var res *http.Response
		var err error

		if test.inMethod == "GET" {
			res, err = http.Get(server.URL + test.inPath)
		}
		if test.inMethod == "POST" {
			res, err = client.PostForm(server.URL+test.inPath, test.inParams)
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

		actualBody, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		if strings.TrimSpace(string(actualBody)) != test.expectedBody {
			t.Errorf("Expected: %s, got: %s", test.expectedBody, actualBody)
		}
	}
}

func genRandomPassword() string {
	b := make([]byte, 10)
	rand.Read(b)
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
}

func TestGracefulShutdown(t *testing.T) {
	server := httptest.NewServer(Handler())
	numRequests := 100
	wg := &sync.WaitGroup{}

	for i := 0; i < numRequests; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			form := url.Values{"password": []string{genRandomPassword()}}
			res, err := http.PostForm(server.URL+"/hash", form)
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
	}

	wg.Add(1)
	time.AfterFunc(1*time.Second, func() {
		defer wg.Done()

		res, err := http.Get(server.URL + "/shutdown")
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != 200 {
			t.Errorf("Expected: %d, got: %d", 200, res.StatusCode)
		}

		res, _ = http.Get(server.URL + "/stats")
		bytes, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(bytes))

		Wait()

		t.Log("Server is shutting down...")
		server.Close()

		res, err = http.Get(server.URL + "/shutdown")
		if err == nil {
			t.Errorf("Expected: %v, got: %v", err, res)
		}
	})

	wg.Wait()
}
