.PHONY: test
test: 
	go test -v -race ./...

.PHONY: install
install: test
	go install ./...

.PHONY: run
run: install
	pwhash