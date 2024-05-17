UNAME := $(shell uname)
RKE_VERSION := $(shell sh -c 'grep -m1 github.com/rancher/rke go.mod | cut -d" " -f2')
SHA ?= $(shell git describe --match=none --always --abbrev=8 --dirty)

IMAGE_TAG ?= dev
REGISTRY ?= docker.io/tlebreux

GO_VERSION ?= 1.22
KUBECTL_VERSION ?= v1.27.10
LOGLEVEL_VERSION ?= v0.1.5
CONTROLLER_GEN_VERSION ?= v0.12.0

# TODO: Make Arch independent
KUBECTL_URL ?= https://storage.googleapis.com/kubernetes-release/release/$(KUBECTL_VERSION)/bin/linux/amd64/kubectl
LOGLEVEL_URL ?= https://github.com/rancher/loglevel/releases/download/$(LOGLEVEL_VERSION)/loglevel-amd64-$(LOGLEVEL_VERSION).tar.gz

IMAGE_TOOLS ?= registry.suse.com/bci/golang:$(GO_VERSION)
IMAGE_MICRO ?= registry.suse.com/bci/bci-micro:15.5
IMAGE_BASE ?= registry.suse.com/bci/bci-base:15.5

LINKFLAGS :=

CGO_ENABLED ?= 0
GO_LDFLAGS ?=
GO_GCFLAGS ?=
GO_BUILDFLAGS ?=

GO_BUILDFLAGS += -tags k8s

WITH_DEBUG ?= false

ifneq (, $(filter $(WITH_DEBUG), t true TRUE y yes 1))
	GO_GCFLAGS += all=-N -l
endif

ifneq ($(UNAME), Darwin)
	GO_LDFLAGS += -extldflags -static
	ifeq (, $(filter $(WITH_DEBUG), t true TRUE y yes 1))
		GO_LDFLAGS += -s
	endif
endif

# FIX: Not working I think..
DEFAULT_VALUES = '{\"rke-version\":\"${RKE_VERSION}\"}'

RANCHER_SERVER_GO_LDFLAGS := -X github.com/rancher/rancher/pkg/version.Version=$(VERSION)
RANCHER_SERVER_GO_LDFLAGS += -X github.com/rancher/rancher/pkg/version.GitCommit=$(VERSION)
RANCHER_SERVER_GO_LDFLAGS += -X github.com/rancher/rancher/pkg/settings.InjectDefaults=$(DEFAULT_VALUES)

RANCHER_AGENT_GO_LDFLAGS := -X main.Version=$(VERSION)

COMMON_ARGS := --file=Dockerfile
COMMON_ARGS += --build-arg=IMAGE_TOOLS=$(IMAGE_TOOLS)
COMMON_ARGS += --build-arg=IMAGE_MICRO=$(IMAGE_MICRO)
COMMON_ARGS += --build-arg=IMAGE_BASE=$(IMAGE_BASE)
COMMON_ARGS += --build-arg=CGO_ENABLED=$(CGO_ENABLED)
COMMON_ARGS += --build-arg=GO_LDFLAGS="$(GO_LDFLAGS)"
COMMON_ARGS += --build-arg=GO_GCFLAGS="$(GO_GCFLAGS)"
COMMON_ARGS += --build-arg=GO_BUILDFLAGS="$(GO_BUILDFLAGS)"
COMMON_ARGS += --build-arg=CONTROLLER_GEN_VERSION="$(CONTROLLER_GEN_VERSION)"
COMMON_ARGS += --build-arg=KUBECTL_URL="$(KUBECTL_URL)"
COMMON_ARGS += --build-arg=LOGLEVEL_URL="$(LOGLEVEL_URL)"

CI_ARGS ?=

BUILD := docker buildx build

TARGETS := $(shell ls scripts)

.dapper:
	@echo Downloading dapper
	@curl -sL https://releases.rancher.com/dapper/latest/dapper-`uname -s`-`uname -m` > .dapper.tmp
	@@chmod +x .dapper.tmp
	@./.dapper.tmp -v
	@mv .dapper.tmp .dapper

$(TARGETS): .dapper
	@if [ "$@" = "check-chart-kdm-source-values" ]; then \
		./.dapper -q --no-out $@; \
	else \
		./.dapper $@; \
	fi

target-%:
	@$(BUILD) \
		--target=$* \
		$(COMMON_ARGS) \
		$(TARGET_ARGS) \
		$(CI_ARGS) .

local-%:
	@$(MAKE) target-$* TARGET_ARGS="--output=type=local,dest=$(DEST) $(TARGET_ARGS)"

registry-%: ## Builds the specified target defined in the Dockerfile using the image/registry output type. The build result will be pushed to the registry if PUSH=true.
	@$(MAKE) target-$* TARGET_ARGS="--output type=image,name=$(REGISTRY)/$*:$(IMAGE_TAG) $(TARGET_ARGS)"

generate:
	@$(MAKE) local-go-generator DEST=./

binary-server:
	@$(MAKE) local-binary-server DEST=./build/server GO_LDFLAGS="$(RANCHER_SERVER_GO_LDFLAGS) $(GO_LDFLAGS)"

binary-agent:
	@$(MAKE) local-binary-agent DEST=./build/agent GO_LDFLAGS="$(RANCHER_AGENT_GO_LDFLAGS) $(GO_LDFLAGS)"

binary: binary-server binary-agent


.DEFAULT_GOAL := ci

.PHONY: $(TARGETS)
