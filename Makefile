# Go parameters
GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
MAIN_PATH=./cmd/etcdtest

# Make parameters
.PHONY: all build test clean run deps ci install-mockgen mockgen build/apiserver build/controller build/kubelet

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

test:
	$(GOTEST) -v ./...

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

# Output directory
OUT_DIR=./out

# Binary names
APISERVER_BINARY=$(OUT_DIR)/apiserver
CONTROLLER_BINARY=$(OUT_DIR)/controller
KUBELET_BINARY=$(OUT_DIR)/kubelet

# Main paths
APISERVER_MAIN=./apiserver/main.go
CONTROLLER_MAIN=./controller/main.go
KUBELET_MAIN=./kubelet/main.go

# Ensure the output directory exists
$(OUT_DIR):
	mkdir -p $(OUT_DIR)

# Build targets
$(OUT_DIR)/%: $(OUT_DIR)
	@$(GOBUILD) -o $(@) -v ./cmd/$(@F)/main.go
	@printf "Built %s\n" $(@F)

build/apiserver: $(APISERVER_BINARY)
build/controller: $(CONTROLLER_BINARY)
build/kubelet: $(KUBELET_BINARY)

# Combined build target
build-all: $(APISERVER_BINARY) $(CONTROLLER_BINARY) $(KUBELET_BINARY)

clean:
	@$(GOCLEAN)
	@rm -f $(APISERVER_BINARY) $(CONTROLLER_BINARY) $(KUBELET_BINARY)
	@rm -rf $(OUT_DIR)
	@printf "Cleaned up build artifacts\n"
