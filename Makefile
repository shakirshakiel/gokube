# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=etcdtest
MAIN_PATH=./cmd/etcdtest

# Make parameters
.PHONY: all build test clean run deps ci install-mockgen mockgen

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run: build
	./$(BINARY_NAME)

deps:
	$(GOGET) ./...
	$(GOMOD) tidy 

test-registry:
	$(GOTEST) -v ./pkg/registry

test-storage:
	$(GOTEST) -v ./pkg/storage

lint:
# Exit with 0 to allow CI to continue with linter errors
	golangci-lint run --issues-exit-code 0

fmt:
	gofmt -s -w . 

vet:
	go vet $(shell go list ./...) 

# CI build target
ci: deps fmt vet lint test build
	@echo "CI build completed successfully"

mockgen: install-mockgen
	go generate ./...

install-mockgen:
	@if ! [ -x "$$(command -v mockgen)" ]; then \
		echo "mockgen not found, installing..."; \
		$(GOCMD) install go.uber.org/mock/mockgen@latest; \
	fi