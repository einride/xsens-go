# Code generated by go.einride.tech/sage. DO NOT EDIT.
# To learn more, see .sage/sagefile.go and https://github.com/einride/sage.

.DEFAULT_GOAL := all

cwd := $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
sagefile := $(abspath $(cwd)/.sage/bin/sagefile)

# Setup Go.
go := $(shell command -v go 2>/dev/null)
export GOWORK ?= off
ifndef go
SAGE_GO_VERSION ?= 1.18.4
export GOROOT := $(abspath $(cwd)/.sage/tools/go/$(SAGE_GO_VERSION)/go)
export PATH := $(PATH):$(GOROOT)/bin
go := $(GOROOT)/bin/go
os := $(shell uname | tr '[:upper:]' '[:lower:]')
arch := $(shell uname -m)
ifeq ($(arch),x86_64)
arch := amd64
endif
$(go):
	$(info installing Go $(SAGE_GO_VERSION)...)
	@mkdir -p $(dir $(GOROOT))
	@curl -sSL https://go.dev/dl/go$(SAGE_GO_VERSION).$(os)-$(arch).tar.gz | tar xz -C $(dir $(GOROOT))
	@touch $(GOROOT)/go.mod
	@chmod +x $(go)
endif

.PHONY: $(sagefile)
$(sagefile): $(go)
	@cd .sage && $(go) mod tidy && $(go) run .

.PHONY: sage
sage:
	@$(MAKE) $(sagefile)

.PHONY: update-sage
update-sage: $(go)
	@cd .sage && $(go) get -d go.einride.tech/sage@latest && $(go) mod tidy && $(go) run .

.PHONY: clean-sage
clean-sage:
	@git clean -fdx .sage/tools .sage/bin .sage/build

.PHONY: all
all: $(sagefile)
	@$(sagefile) All

.PHONY: convco-check
convco-check: $(sagefile)
	@$(sagefile) ConvcoCheck

.PHONY: format-markdown
format-markdown: $(sagefile)
	@$(sagefile) FormatMarkdown

.PHONY: format-yaml
format-yaml: $(sagefile)
	@$(sagefile) FormatYAML

.PHONY: git-verify-no-diff
git-verify-no-diff: $(sagefile)
	@$(sagefile) GitVerifyNoDiff

.PHONY: go-generate
go-generate: $(sagefile)
	@$(sagefile) GoGenerate

.PHONY: go-lint
go-lint: $(sagefile)
	@$(sagefile) GoLint

.PHONY: go-mod-tidy
go-mod-tidy: $(sagefile)
	@$(sagefile) GoModTidy

.PHONY: go-review
go-review: $(sagefile)
	@$(sagefile) GoReview

.PHONY: go-test
go-test: $(sagefile)
	@$(sagefile) GoTest

.PHONY: stringer
stringer: $(sagefile)
	@$(sagefile) Stringer
