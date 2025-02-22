DIST_DIR=dist
DIST_TARGETS=$(addprefix dist/,$(BINARIES))

install-goreleaser: ## Install goreleaser
	@if ! [ -x "$$(command -v goreleaser)" ]; then \
		echo "goreleaser not found, installing..."; \
		$(GOCMD) install github.com/goreleaser/goreleaser/v2@latest; \
	fi

goreleaser-check: install-goreleaser ## Check gorelease configuration is correct
	@goreleaser check	

$(DIST_DIR): install-goreleaser goreleaser-check ## Build binaries for goreleaser
	@goreleaser build --snapshot --clean

$(DIST_TARGETS): install-goreleaser goreleaser-check ## Build specific binaries for goreleaser
	@goreleaser build --snapshot --clean --id $(@F)

GO_KUBE_RELEASE_BINARIES = $(foreach binary,$(BINARIES),$(HOME)/gokube/$(binary))

$(HOME)/gokube: ## Create gokube directory
	@if [ ! -d $(HOME)/gokube ]; then mkdir -p $(HOME)/gokube; fi

$(GO_KUBE_RELEASE_BINARIES): $(HOME)/gokube ## Copy binaries to gokube
	@echo $(@F) $(basename $(@F))
	@cp $(DIST_DIR)/$(@F)_linux_arm64_v8.0/$(@F) $(HOME)/gokube/
	@printf "Copied linux arm64 binary to $(HOME)/gokube\n"

install-dist: dist $(HOME)/gokube/apiserver $(HOME)/gokube/controller $(HOME)/gokube/kubelet $(HOME)/gokube/scheduler ## Create distributions and copy to gokube directory

# Lima commands for VMs
LIMA_VMS = master worker1 worker2
LIMA_INIT_TARGETS = $(addprefix init/,$(LIMA_VMS))
LIMA_START_TARGETS = $(addprefix start/,$(LIMA_VMS))
LIMA_STOP_TARGETS = $(addprefix stop/,$(LIMA_VMS))
LIMA_DELETE_TARGETS = $(addprefix delete/,$(LIMA_VMS))
LIMA_SHELL_TARGETS = $(addprefix shell/,$(LIMA_VMS))

$(LIMA_START_TARGETS): $(GO_KUBE_RELEASE_BINARIES) ## Start Lima VM
	@limactl start --name=$(@F) workbench/debian-12.yaml --tty=false
	@printf "Lima instance '$(@F)' started\n"

$(LIMA_STOP_TARGETS): ## Stop Lima VM
	@limactl stop $(@F)
	@printf "Lima instance '$(@F)' stopped\n"

$(LIMA_DELETE_TARGETS): ## Delete Lima VM
	@limactl delete $(@F)
	@printf "Lima instance '$(@F)' deleted\n"

$(LIMA_SHELL_TARGETS): ## Go to shell of Lima VM
	@printf "Entering Lima instance '$(@F)' shell\n"
	@limactl shell --workdir $(HOME) $(@F)

lima/init-vms: $(GO_KUBE_RELEASE_BINARIES) start/master stop/master start/worker1 stop/worker1 start/worker2 stop/worker2 ## Init Lima VMs

lima/start-vms: start/master start/worker1 start/worker2 ## Start all Lima VMs

lima/run: ### Run the project
	process-compose -f process-compose-lima.yml up

lima/cleanup: stop/master stop/worker1 stop/worker2 delete/master delete/worker1 delete/worker2 # Cleanup all lima vms