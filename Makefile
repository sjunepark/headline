.PHONY: build

check:
	go fmt ./...
	golangci-lint run
	go vet ./...

test:
	go test ./...

build:
	wails build