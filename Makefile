.PHONY: run build test lint

run:
	go run main.go

build:
	go build -o ditto-cli main.go

test:
	go test ./...

lint:
	golangci-lint run
