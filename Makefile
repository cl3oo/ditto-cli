.PHONY: run build test lint smoke clean

GOBIN ?= $(CURDIR)/.bin
DIST_DIR ?= $(CURDIR)/dist
APP_NAME ?= ditto-cli
APP_BIN := $(DIST_DIR)/$(APP_NAME)
GOLANGCI_LINT := $(GOBIN)/golangci-lint
GOLANGCI_LINT_PKG := github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

run:
	go run . tui

build:
	@mkdir -p $(DIST_DIR)
	go build -o $(APP_BIN) .

test:
	go test ./...

$(GOLANGCI_LINT):
	@mkdir -p $(GOBIN)
	GOBIN=$(GOBIN) go install $(GOLANGCI_LINT_PKG)

lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run

smoke: build
	./scripts/smoke.sh

clean:
	rm -rf $(DIST_DIR)
