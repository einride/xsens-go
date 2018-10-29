.PHONY: all
all: lint test-go

include build/rules.mk
build/rules.mk:
	git submodule update --init

lint: lint-vendor lint-go lint-markdown

.PHONY: test-go
test-go: vendor
	go test -race -cover ./...

.PHONY: doc
doc:
	godoc -http=:6060 &
	xdg-open http://localhost:6060/pkg/github.com/einride/xsens-go

