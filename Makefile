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

# Colors
CYAN_COLOR_START := \033[36m
CYAN_COLOR_END := \033[0m

.PHONY: all build test clean run deps ci install-mockgen mockgen $(BUILD_TARGETS) $(DIST_TARGETS) $(INSTALL_TARGETS) $(GO_BIN_TARGETS)

help: ## Prints help (only for targets with comments)
	@grep -E '^[a-zA-Z._/\-]+:.*?## ' $(MAKEFILE_LIST) | sort | awk -F'[:##]' '{printf "$(CYAN_COLOR_START)%-15s $(CYAN_COLOR_END)%s\n", $$1, $$4}'

deps: ## Install/Upgrade dependencies
	$(GOGET) ./...
	$(GOMOD) tidy

fmt: ## Format go code
	gofmt -s -w .

vet: ## Run SCA using go vet
	go vet $(shell go list ./...)

lint: ## Run lint
	docker run --rm -v $(PWD):/app -v $(PWD)/.golangci-lint-cache:/root/.cache -w /app golangci/golangci-lint:v1.63.4 golangci-lint run -v --exclude S1000

test: ## Run all tests
	$(GOTEST) -v ./...

test/%: ## Run package level tests
	$(GOTEST) -v ./pkg/$(@F)

ci: deps fmt vet lint test ## Run CI test target(deps,fmt,vet,lint,test)
	@echo "CI build completed successfully"

mockgen: install-mockgen
	go generate ./...

install-mockgen:
	@if ! [ -x "$$(command -v mockgen)" ]; then \
		echo "mockgen not found, installing..."; \
		$(GOCMD) install go.uber.org/mock/mockgen@latest; \
	fi

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

# Output directory
$(DIST_DIR):
	@goreleaser build --snapshot --clean

$(DIST_TARGETS):
	@goreleaser build --snapshot --clean --id $(@F)

GO_KUBE_RELEASE_BINARIES = $(foreach binary,$(BINARIES),$(HOME)/gokube/$(binary))

$(HOME)/gokube:
	@mkdir -p $(HOME)/gokube

$(GO_KUBE_RELEASE_BINARIES): $(HOME)/gokube
	@echo $(@F) $(basename $(@F))
	@cp $(DIST_DIR)/$(@F)_linux_arm64/$(@F) $(HOME)/gokube
	@printf "Copied linux arm64 binary to $(HOME)/gokube\n"

clean:
	@$(GOCLEAN)
	@rm -f $(BINARY_PATHS)
	@rm -rf $(OUT_DIR)
	@printf "Cleaned up build artifacts\n"
	@rm -f $(EXECUTABLES)
	@printf "Cleaned up installed binaries\n"
	@rm -rf $(DIST_DIR)
	@printf "Cleaned up dist artifacts\n"
	@rm -rf $(HOME)/gokube
	@printf "Cleaned up gokube binaries\n"

# Lima commands for VMs
LIMA_VMS = master worker1
LIMA_START_TARGETS = $(addprefix start/,$(LIMA_VMS))
LIMA_STOP_TARGETS = $(addprefix stop/,$(LIMA_VMS))
LIMA_DELETE_TARGETS = $(addprefix delete/,$(LIMA_VMS))
LIMA_SHELL_TARGETS = $(addprefix shell/,$(LIMA_VMS))
LIMA_TARGETS = $(LIMA_START_TARGETS) $(LIMA_STOP_TARGETS) $(LIMA_DELETE_TARGETS) $(LIMA_SHELL_TARGETS)

$(LIMA_START_TARGETS): $(GO_KUBE_RELEASE_BINARIES)
	@limactl start --name=$(@F) workbench/debian-12.yaml --tty=false
	@printf "Lima instance '$(@F)' started\n"

$(LIMA_STOP_TARGETS):
	@limactl stop $(@F)
	@printf "Lima instance '$(@F)' stopped\n"

$(LIMA_DELETE_TARGETS):
	@limactl delete $(@F)
	@printf "Lima instance '$(@F)' deleted\n"

$(LIMA_SHELL_TARGETS):
	@printf "Entering Lima instance '$(@F)' shell\n"
	@limactl shell --workdir $(HOME) $(@F)
