.PHONY: all
all: markdown-lint go-lint go-test

include build/rules.mk
build/rules.mk:
	git submodule update --init --recursive

.PHONY: go-test
go-test: dep-ensure
	go test -race -cover ./...

.PHONY: go-lint
go-lint: $(GOLANGCI_LINT) dep-ensure
	$(GOLANGCI_LINT) run --enable-all --skip-dirs vendor

.PHONY: dep-ensure
dep-ensure: $(DEP)
	$(DEP) ensure -v

.PHONY: markdown-lint
markdown-lint: $(MARKDOWNLINT)
	$(MARKDOWNLINT) --ignore build --ignore vendor .

.PHONY: doc
doc:
	godoc -http=:6060 &
	xdg-open http://localhost:6060/pkg/github.com/einride/xsens-go
