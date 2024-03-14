#!/usr/bin/make -f

export VERSION := $(shell echo $(shell git describe --always --match "v*") | sed 's/^v//')
export TMVERSION := $(shell go list -m github.com/cometbft/cometbft | sed 's:.* ::')
export COMMIT := $(shell git log -1 --format='%H')
export MAKE_PROJECT_ROOT := $(CURDIR)
LEDGER_ENABLED ?= true
BINDIR ?= $(GOPATH)/bin
BUILDDIR ?= $(CURDIR)/build
DOCKER := $(shell which docker)
PROJECT_NAME = $(shell git remote get-url origin | xargs basename -s .git)
COSMOS_VERSION = $(shell go list -m github.com/cosmos/cosmos-sdk | sed 's:.* ::')
IGNITE_VERSION = v0.27.1
MOCKS_DIR = $(CURDIR)/tests/mocks

# $(info GOOS: $(GOOS), GOARCH: $(GOARCH), CC: $(CC), CXX: $(CXX))

# process build tags
build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
	ifeq ($(OS),Windows_NT)
	GCCEXE = $(shell where gcc.exe 2> NUL)
	ifeq ($(GCCEXE),)
		$(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
	else
		build_tags += ledger
	endif
	else
	UNAME_S = $(shell uname -s)
	ifeq ($(UNAME_S),OpenBSD)
		$(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
	else
		GCC = $(shell command -v gcc 2> /dev/null)
		ifeq ($(GCC),)
			$(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
		else
			build_tags += ledger
		endif
	endif
	endif
endif

# this is required by ethermint
build_tags += cgo

ifeq (secp,$(findstring secp,$(COSMOS_BUILD_OPTIONS)))
  build_tags += libsecp256k1_sdk
endif
# DB backend selection
ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += gcc
endif
ifeq (badgerdb,$(findstring badgerdb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += badgerdb
endif
# handle rocksdb
ifeq (rocksdb,$(findstring rocksdb,$(COSMOS_BUILD_OPTIONS)))  
  build_tags += rocksdb
endif
# handle boltdb
ifeq (boltdb,$(findstring boltdb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += boltdb
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))
export MAKE_BUILD_TAGS := $(build_tags)
export MAKE_BUILD_TAGS_COMMA_SEP := $(build_tags_comma_sep)

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=titan \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=titand \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
		  -X github.com/cometbft/cometbft/version.TMCoreSemVer=$(TMVERSION)

ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ifeq (static,$(findstring static,$(COSMOS_BUILD_OPTIONS)))
	ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
else
	ldflags += -extldflags "-Wl,-rpath,$$ORIGIN/../lib -Wl,-rpath,@executable_path/../lib -Wl,-rpath,$$ORIGIN/lib -Wl,-rpath,@executable_path/lib -Wl,-rpath,$$ORIGIN -Wl,-rpath,@executable_path"
endif

ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
# check for nostrip option
ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

# Check for debug option
ifeq (debug,$(findstring debug,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -gcflags "all=-N -l"
endif

all: build-with-regen

###############################################################################
###                                Linting                                  ###
###############################################################################

lint:	golangci-lint
	go mod verify	
	golangci-lint run --out-format=tab --timeout 2m0s

lint-fix:	golangci-lint	
	golangci-lint run --fix --out-format=tab --issues-exit-code=0 --timeout 2m0s

.PHONY: lint lint-fix

format:
	find . -name '*.go' -type f -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -name '*.pb.go' -not -name '*.pb.gw.go' | xargs gofumpt -w -l

.PHONY: format

###############################################################################
###                                  Build                                  ###
###############################################################################

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILDDIR)/

$(BUILD_TARGETS): go.sum $(BUILDDIR)/
	CGO_ENABLED=1 go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./...

build-with-regen: proto-all lint go.sum $(BUILDDIR)/
	go build -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./...

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

.PHONY: build build-with-regen build-linux-amd64 build-linux-arm64 

clean:
	rm -rf $(BUILDDIR)/

.PHONY: clean

mocks: $(MOCKS_DIR)
	@go install github.com/golang/mock/mockgen@v1.6.0
	sh ./scripts/mockgen.sh
.PHONY: mocks

$(MOCKS_DIR):
	mkdir -p $(MOCKS_DIR)

###############################################################################
###                                Protobuf                                 ###
###############################################################################

protoVer=0.11.2
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoContainerName=protobuilder-titan-$(subst /,_,$(subst \,_,$(CURDIR)))
protoImage=$(DOCKER) run -v $(CURDIR):/workspace --workdir /workspace --user 0 --name $(protoContainerName) $(protoImageName)
protoFormatImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace --user 0 $(protoImageName)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"	
	@if [ ! "$(shell docker ps -aq -f name=$(protoContainerName))" ]; then \
		$(protoImage) sh ./scripts/protocgen.sh ; \
	else \
		$(DOCKER) start -a $(protoContainerName) ; \
	fi
	

proto-format:
	@$(protoFormatImage) find ./ -name "*.proto" -exec clang-format --style="{ IndentWidth: 2, BasedOnStyle: google, AlignConsecutiveAssignments: true, AlignConsecutiveDeclarations: Consecutive, ColumnLimit: 120 }" -i {} \;

proto-lint:
	@$(protoFormatImage) buf lint --error-format=json


###############################################################################
###                              Documentation                              ###
###############################################################################

update-swagger-docs-by-ignite:
	@ignite generate openapi -y

.PHONY: update-swagger-docs-by-ignite

###############################################################################
###                              Dev tools		                              ###
###############################################################################

ignite:
	@echo "Installing ignite (tag ${IGNITE_VERSION}) ..."
	rm -rf ignite_tmp	
	git clone --depth 1 --branch ${IGNITE_VERSION} https://github.com/ignite/cli.git ignite_tmp	
	cd ignite_tmp && make install
	rm -rf ignite_tmp

cosmovisor:
	@echo "Installing cosmovisor ..."
	rm -rf cosmovisor_tmp	
	git clone --depth 1 --branch ${COSMOS_VERSION} https://github.com/cosmos/cosmos-sdk.git cosmovisor_tmp
	cd cosmovisor_tmp && make cosmovisor 	
	cp cosmovisor_tmp/tools/cosmovisor/cosmovisor build/cosmovisor
	rm -rf cosmovisor_tmp

GOLANGCI_VERSION=v1.52.2

golangci-lint:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_VERSION)

.PHONE: ignite cosmovisor golangci-lint
###############################################################################
###                                Localnet                                 ###
###############################################################################

localnet-serve:
# TODO : because ignite gen proto-go have different result with make proto-gen => so skip-proto
# 			but we must findout why and sync 2 way to generate proto-go
	@ignite chain serve --skip-proto	-v


localnet-serve-reset:
	@ignite chain serve --skip-proto --reset-once -v

###############################################################################
###                                Testing                                  ###
###############################################################################

test-solidity:
	@echo "Beginning solidity tests..."
	./scripts/run-solidity-tests.sh

vulncheck: $(BUILDDIR)/
	GOBIN=$(BUILDDIR) go install golang.org/x/vuln/cmd/govulncheck@latest
	$(BUILDDIR)/govulncheck ./...

go.sum: go.mod
	@echo "Ensure dependencies have not been modified ..." >&2
	go mod verify
	go mod tidy

test-testutil:
	go test -timeout 1200s -cover github.com/titantkx/titan/testutil -v

test-unit:
	go test -v -timeout 1200s -cover github.com/titantkx/titan/x/...

test-app:
	go test -timeout 1200s -cover github.com/titantkx/titan/app -v

test-integration:	
	go test -timeout 1200s -cover github.com/titantkx/titan/tests/integration/... -v

test-e2e-cmd: 
	TEST_TYPE=basic go test -timeout 1200s -count=1 github.com/titantkx/titan/tests/e2e/cmd -v

test-all: test-testutil test-unit test-app test-integration test-e2e-cmd

###############################################################################
###                                Releasing                                ###
###############################################################################

PACKAGE_NAME:=github.com/titantkx/titan
GOLANG_CROSS_VERSION  = v1.20
GOPATH ?= '$(HOME)/go'
release-dry-run:
	./scripts/release.sh --dry-run

release:
	./scripts/release.sh

release-gen-info:
	@./scripts/gen-upgrade-info.sh v$(VERSION)

.PHONY: release-dry-run release release-gen-info