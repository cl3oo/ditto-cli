.PHONY: run build test lint

run:
	go run . tui

build:
	go build -o ditto-cli main.go

test:
	go test ./...

lint:
	golangci-lint run
