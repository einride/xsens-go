SHELL := /bin/bash

.PHONY: all
all: \
	commitlint \
	prettier-markdown \
	mockgen-generate \
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

.PHONY: go-test
go-test:
	$(info [$@] running Go tests...)
	@go test -count 1 -cover -race ./...

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
	pkg/serial/baudrate_string.go

coordinatesystem_string.go: coordinatesystem.go $(stringer)
	$(info generating $@...)
	@$(stringer) -type CoordinateSystem -trimprefix CoordinateSystem -output $@ $<

datatype_string.go: datatype.go $(stringer)
	$(info generating $@...)
	@$(stringer) -type DataType -trimprefix DataType -output $@ $<

errorcode_string.go: errorcode.go $(stringer)
	$(info generating $@...)
	@$(stringer) -type ErrorCode -trimprefix ErrorCode -output $@ $<

fixtype_string.go: fixtype.go $(stringer)
	$(info generating $@...)
	@$(stringer) -type FixType -trimprefix FixType -output $@ $<

messageidentifier_string.go: messageidentifier.go $(stringer)
	$(info generating $@...)
	@$(stringer) -type MessageIdentifier -trimprefix MessageIdentifier -output $@ $<

precision_string.go: precision.go $(stringer)
	$(info generating $@...)
	@$(stringer) -type Precision -trimprefix Precision -output $@ $<

pkg/serial/baudrate_string.go: pkg/serial/baudrate.go $(stringer)
	$(info generating $@...)
	@$(stringer) -type BaudRate -trimprefix BaudRate -output $@ $<

.PHONY: mockgen-generate
mockgen-generate: test/mocks/xsens/mocks.go

test/mocks/xsens/mocks.go: client.go go.mod
	go run github.com/golang/mock/mockgen \
		-destination $@ -package mockxsens \
		github.com/einride/xsens-go SerialPort
