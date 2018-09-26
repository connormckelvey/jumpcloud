GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=mybinary
BINARY_UNIX=$(BINARY_NAME)_unix

.PHONY: test
test: 
	go test -v -race ./...

.PHONY: install
install: test
	go install ./...

.PHONY: run
run: install
	pwhash