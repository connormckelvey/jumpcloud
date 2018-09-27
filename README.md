# JumpCloud 

**PWHASH** - The Password Hashing Service

## Installing

- `go get github.com/connormckelvey/jumpcloud/cmd/pwhash`

## Running

- `pwhash` (Start with defaults)
- `pwhash -h`

```
Usage of pwhash:
  -d int
    	Amount to delay hash response in seconds (default 5)
  -p int
    	Port to listen on (default 8080)
```

## Testing

- `cd $GOPATH/src/github.com/connormckelvey/jumpcloud`
- `make test`


## Assumptions and Design

### Stats / Metrics

The `/stats` endpoint should be available as long as possible. When a shutdown has started, all requests to other endpoints like `/hash` and `/shutdown` immediately return 503 responses, but `/stats` will continue to successfully respond to requests until the `http.Server` is shutdown. This way, metrics can continue to be collected until all `/hash` requests have completed. `/stats` can also be used as a health check endpoint for a system like Kubernetes, preventing an unwanted shutdown.

### Shutdown

`http.Server.Shutdown` with a timeout Context isn't good enough. Even though open connections are allowed to complete within the timeout, a long running request may get killed when the timeout is reached. Using a `sync.WaitGroup`, the application waits for all `/hash` requests to complete before calling `http.Server.Shutdown`

### External Dependencies

Because this was a coding challenge, I chose to not use any 3rd party libraries. I believe the packages I have created are of "production quality," but really, I would choose to use something like Logrus for logging and Prometheus for metrics.

### HTTP Servers == Concurrency 

All requests to an `http.Server` are run in separate goroutines, so ensuring that any
shared state is accessed in a safe way is important in preventing data races and in preventing corruption of that state. Using things like `sync.Mutex`, channels, and atomic updates makes concurrent access safe.

### Dependency Injection

Throughout the service, handlers and middleware may need access to common things like the metrics store, configuration information, and the logger. By defining handlers and middleware as methods on  `Application`, injecting a dependency means simply adding that dependency to `Application`.

## What I didn't quite get to

### Hide passwords in logs

There is no logging in the `/hash` handler, but if this were a production system, it would be pretty easy to accidentally log a password. I had plans to create an `io.Writer` and update the `withLogging` middleware to scan the bytes being written and obfuscate the password before writing to StdOut.

### TLS

Im many cases, a microservice may not need to be the termination point for a TLS connection and a TLS sidecar / reverse proxy will do. But in this service, the main functionality is password hashing, so it would make the most sense to actually start a TLS server to provide end-to-end encryption.
