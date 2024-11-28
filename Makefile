# Go parameters
GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOINSTALL=$(GOCMD) install
MAIN_PATH=./cmd/etcdtest

# Make parameters
OUT_DIR=out
DIST_DIR=dist
BINARIES=apiserver controller kubelet
BINARY_PATHS=$(addprefix $(OUT_DIR)/,$(BINARIES))
EXECUTABLES=$(addprefix $(GOPATH)/,$(BINARIES))

BUILD_TARGETS=$(addprefix build/,$(BINARIES))
DIST_TARGETS=$(addprefix dist/,$(BINARIES))
INSTALL_TARGETS=$(addprefix install/,$(BINARIES))
GO_BIN_TARGETS=$(addprefix $(GOPATH)/bin/,$(BINARIES))

.PHONY: all build test clean run deps ci install-mockgen mockgen $(DIST_DIR) $(BUILD_TARGETS) $(DIST_TARGETS) $(INSTALL_TARGETS) $(GO_BIN_TARGETS)

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
$(DIST_DIR):
	@goreleaser build --snapshot --clean

$(DIST_TARGETS):
	@goreleaser build --snapshot --clean --id $(@F)

# Main paths
# Ensure the output directory exists
$(OUT_DIR):
	@mkdir -p $(OUT_DIR)

# Build targets
$(OUT_DIR)/%: $(OUT_DIR)
	@$(GOBUILD) -o $(@) -v ./cmd/$(@F)/$(@F).go
	@printf "Built %s\n" $(@F)

build/apiserver: $(OUT_DIR)/apiserver
build/controller: $(OUT_DIR)/controller
build/kubelet: $(OUT_DIR)/kubelet

$(GO_BIN_TARGETS):
	@printf "Installing %s...\n" $(@F)
	@$(GOINSTALL) ./cmd/$(@F)/$(@F).go
	@printf "Successfully installed %s\n" $(@F)
	@printf "Executable located at %s\n\n" $(GOPATH)/bin/$(@F)

install/apiserver: $(GOPATH)/bin/apiserver
install/controller: $(GOPATH)/bin/controller
install/kubelet: $(GOPATH)/bin/kubelet

# Combined build target
build-all: $(BINARY_PATHS)
install-all: $(EXECUTABLES)

clean:
	@$(GOCLEAN)
	@rm -f $(BINARY_PATHS)
	@rm -rf $(OUT_DIR)
	@printf "Cleaned up build artifacts\n"
	@rm -f $(EXECUTABLES)
	@printf "Cleaned up installed binaries\n"
