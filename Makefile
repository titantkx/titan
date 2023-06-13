#!/usr/bin/make -f

export VERSION := $(shell echo $(shell git describe --always --match "v*") | sed 's/^v//')
export TMVERSION := $(shell go list -m github.com/cometbft/cometbft | sed 's:.* ::')
export COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true
BINDIR ?= $(GOPATH)/bin
BUILDDIR ?= $(CURDIR)/build
DOCKER := $(shell which docker)
PROJECT_NAME = $(shell git remote get-url origin | xargs basename -s .git)

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

ifeq (secp,$(findstring secp,$(COSMOS_BUILD_OPTIONS)))
  build_tags += libsecp256k1_sdk
endif

ifeq (legacy,$(findstring legacy,$(COSMOS_BUILD_OPTIONS)))
  build_tags += app_v1
endif

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=titan \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=titand \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
		  -X github.com/cometbft/cometbft/version.TMCoreSemVer=$(TMVERSION)


# DB backend selection
ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += gcc
endif
ifeq (badgerdb,$(findstring badgerdb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += badgerdb
endif
# handle rocksdb
ifeq (rocksdb,$(findstring rocksdb,$(COSMOS_BUILD_OPTIONS)))
  CGO_ENABLED=1
  build_tags += rocksdb
endif
# handle boltdb
ifeq (boltdb,$(findstring boltdb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += boltdb
endif

ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
# check for nostrip option
ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

# Check for debug option
ifeq (debug,$(findstring debug,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -gcflags "all=-N -l"
endif

all: build

###############################################################################
###                                  Build                                  ###
###############################################################################

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILDDIR)/

build-linux-amd64:
	GOOS=linux GOARCH=amd64 LEDGER_ENABLED=false $(MAKE) build

build-linux-arm64:
	GOOS=linux GOARCH=arm64 LEDGER_ENABLED=false $(MAKE) build

$(BUILD_TARGETS): proto-all go.sum $(BUILDDIR)/
	go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./...

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

.PHONY: build build-linux-amd64 build-linux-arm64 

clean:
	rm -rf $(BUILDDIR)/

.PHONY: clean

###############################################################################
###                                Protobuf                                 ###
###############################################################################

protoVer=0.11.2
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoContainerName=protobuilder-titan-$(subst /,_,$(subst \,_,$(CURDIR)))
protoImage=$(DOCKER) run -v $(CURDIR):/workspace --workdir /workspace --name $(protoContainerName) $(protoImageName)
protoFormatImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	@if [ ! "$(shell docker ps -aq -f name=$(protoContainerName))" ]; then \
    $(protoImage) sh ./scripts/protocgen.sh ; \
	else \
		$(DOCKER) start -a $(protoContainerName) ; \
	fi
	

proto-format:
	@$(protoFormatImage) find ./ -name "*.proto" -exec clang-format --style="{ IndentWidth: 2, BasedOnStyle: google, AlignConsecutiveAssignments: true, AlignConsecutiveDeclarations: Consecutive }" -i {} \;

proto-lint:
	@$(protoFormatImage) buf lint --error-format=json

proto-swagger-gen:
	@echo "Generating Protobuf Swagger"
	@$(protoFormatImage) sh ./scripts/protoc-swagger-gen.sh	
	$(MAKE) update-swagger-docs

###############################################################################
###                              Documentation                              ###
###############################################################################

statik: 
	@echo "Installing statik..."	
	@go install github.com/rakyll/statik@v0.1.6

update-swagger-docs: statik
	statik -src=client/docs/swagger-ui -dest=client/docs -f -m
	@if [ -n "$(git status --porcelain)" ]; then \
		echo "\033[91mSwagger docs are out of sync!!!\033[0m";\
		exit 1;\
	else \
		echo "\033[92mSwagger docs are in sync\033[0m";\
	fi
.PHONY: update-swagger-docs

update-swagger-docs-by-ignite:
	@ignite generate openapi -y


###############################################################################
###                              Dev tools		                              ###
###############################################################################

ignite:
	@echo "Installing ignite (tag nightly) ..."
	rm -rf ignite_tmp	
	git clone --depth 1 --branch nightly https://github.com/ignite/cli.git ignite_tmp	
	cd ignite_tmp && make install
	rm -rf ignite_tmp


###############################################################################
###                                Localnet                                 ###
###############################################################################

localnet-serve:
# TODO : because ignite gen proto-go have different result with make proto-gen => so skip-proto
# 			but we must findout why and sync 2 way to generate proto-go
	@ignite chain serve --skip-proto