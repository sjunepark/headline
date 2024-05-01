.PHONY: build

check:
	go fmt ./...
	golangci-lint run
	go vet ./...

check-ci:
	go fmt ./...
	go vet ./...

test:
	go test ./...

coverage:
	go test -coverprofile=c.out ./...
	go tool cover -html="c.out"

build:
	wails build