# all: run all targets included in a full build
.PHONY: all
all: \
	markdown-lint \
	go-generate \
	go-lint \
	go-test \
	go-mod-tidy \
	git-verify-nodiff \
	git-verify-submodules

export GO111MODULE := on

.PHONY: build
build:
	@git submodule update --init --recursive $@

include build/rules.mk
build/rules.mk: build
	@# included in submodule: build

# go-test: run Go test suite
.PHONY: go-test
go-test:
	go test -race -cover ./...

# go-lint: lint Go code
.PHONY: go-lint
go-lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run ./... --enable-all --skip-dirs vendor

# go-mod-tidy: ensure Go module files are in sync
.PHONY: go-mod-tidy
mod-tidy:
	go mod tidy -v

# markdown-lint: lint Markdown documentation
.PHONY: markdown-lint
markdown-lint: $(MARKDOWNLINT)
	$(MARKDOWNLINT) --ignore build --ignore vendor .

# go-generate: run Go code generators
.PHONY: go-generate
go-generate: $(GOBIN)
	go generate ./...
