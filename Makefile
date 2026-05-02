.PHONY: run build test lint smoke clean version release-snapshot release-check release-build release-test

GOBIN ?= $(CURDIR)/.bin
DIST_DIR ?= $(CURDIR)/dist
APP_NAME ?= ditto-cli
APP_BIN := $(DIST_DIR)/$(APP_NAME)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w -X github.com/rfcku/ditto-cli/cmd.version=$(VERSION) -X github.com/rfcku/ditto-cli/cmd.commit=$(COMMIT) -X github.com/rfcku/ditto-cli/cmd.buildDate=$(BUILD_DATE)
BUILD_FLAGS := -ldflags "$(LDFLAGS)"
REVIVE := $(GOBIN)/revive
REVIVE_PKG := github.com/mgechev/revive@latest

run:
	go run $(BUILD_FLAGS) . tui

build:
	@mkdir -p $(DIST_DIR)
	go build $(BUILD_FLAGS) -o $(APP_BIN) .

test:
	go test ./...

$(REVIVE):
	@mkdir -p $(GOBIN)
	GOBIN=$(GOBIN) go install $(REVIVE_PKG)

lint: $(REVIVE)
	$(REVIVE) -config revive.toml -formatter friendly ./...

smoke: build
	./scripts/smoke.sh

clean:
	rm -rf $(DIST_DIR)

version:
	@echo $(VERSION)

release-snapshot:
	goreleaser release --snapshot --clean

release-check:
	goreleaser check

release-build:
	goreleaser build --snapshot --clean

release-test:
	goreleaser release --snapshot --skip=publish --clean
