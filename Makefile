SHELL := /bin/bash

.PHONY: all
all: \
	commitlint \
	prettier-markdown \
	go-generate \
	go-lint \
	go-review \
	go-test \
	go-mod-tidy \
	git-verify-nodiff

include tools/commitlint/rules.mk
include tools/git-verify-nodiff/rules.mk
include tools/golangci-lint/rules.mk
include tools/goreview/rules.mk
include tools/prettier/rules.mk
include tools/semantic-release/rules.mk
include tools/stringer/rules.mk

.PHONY: clean
clean:
	$(info [$@] cleaning generated files...)
	@find -name '*_string.go' -exec rm {} \+
	@rm -rf internal/gen
	@rm -rf build

.PHONY: go-test
go-test:
	$(info [$@] running Go tests...)
	@mkdir -p build/coverage
	@go test -count 1 -short -race -coverprofile=build/coverage/$@.txt -covermode=atomic ./...

.PHONY: go-mod-tidy
go-mod-tidy:
	$(info [$@] tidying Go module files...)
	@go mod tidy -v

.PHONY: go-generate
go-generate: \
	coordinatesystem_string.go \
	datatype_string.go \
	errorcode_string.go \
	fixtype_string.go \
	messageidentifier_string.go \
	precision_string.go \
	serial/baudrate_string.go \
	internal/gen/mockxsens/mocks.go

%_string.go: %.go $(stringer)
	$(info generating $*.go)
	go generate ./$<

internal/gen/mockxsens/mocks.go: client.go go.mod
	$(info generating $@...)
	@go generate ./$<
