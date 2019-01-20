# all: run all targets included in a full build
.PHONY: all
all: \
	markdown-lint \
	dep-ensure \
	go-lint \
	go-test \
	go-generate \
	git-verify-no-diff

.PHONY: build
build:
	@git submodule update --init --recursive $@

include build/rules.mk
build/rules.mk: build
	@# included in submodule: build

# git-verify-no-diff: verify that the working tree does not contain a diff
.PHONY: git-verify-no-diff
git-verify-no-diff:
	git diff --quiet

# go-test: run Go test suite
.PHONY: go-test
go-test:
	go test -race -cover ./...

# go-lint: lint Go code
.PHONY: go-lint
go-lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run ./... --enable-all --skip-dirs vendor

# dep-ensure: update Go dependencies
.PHONY: dep-ensure
dep-ensure: $(DEP)
	$(DEP) ensure -v

# markdown-lint: lint Markdown documentation
.PHONY: markdown-lint
markdown-lint: $(MARKDOWNLINT)
	$(MARKDOWNLINT) --ignore build --ignore vendor .

# go-generate: run Go code generators
.PHONY: go-generate
go-generate: $(STRINGER)
	go generate ./...
