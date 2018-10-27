include build/rules.mk
build/rules.mk:
	git submodule update --init

lint: lint-vendor lint-go lint-markdown

.PHONY: test-go
test-go: vendor
	go test -race -cover ./...
