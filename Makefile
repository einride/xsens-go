.PHONY: all
all: \
	mockgen-generate \
	go-generate \
	go-review \
	go-lint \
	go-test \
	go-mod-tidy \
	git-verify-nodiff

include tools/git-verify-nodiff/rules.mk
include tools/golangci-lint/rules.mk
include tools/goreview/rules.mk

.PHONY: go-test
go-test:
	go test -race -cover ./...

.PHONY: go-mod-tidy
go-mod-tidy:
	find . -name go.mod -execdir go mod tidy -v \;

# go-generate: generate Go code
.PHONY: go-generate
go-generate: \
	coordinatesystem_string.go \
	datatype_string.go \
	errorcode_string.go \
	fixtype_string.go \
	messageidentifier_string.go \
	precision_string.go \
	pkg/serial/baudrate_string.go

stringer := go run -modfile tools/xtools/go.mod golang.org/x/tools/cmd/stringer

coordinatesystem_string.go: coordinatesystem.go tools/xtools/go.mod
	$(info generating $@...)
	@$(stringer) -type CoordinateSystem -trimprefix CoordinateSystem -output $@ $<

datatype_string.go: datatype.go tools/xtools/go.mod
	$(info generating $@...)
	@$(stringer) -type DataType -trimprefix DataType -output $@ $<

errorcode_string.go: errorcode.go tools/xtools/go.mod
	$(info generating $@...)
	@$(stringer) -type ErrorCode -trimprefix ErrorCode -output $@ $<

fixtype_string.go: fixtype.go tools/xtools/go.mod
	$(info generating $@...)
	@$(stringer) -type FixType -trimprefix FixType -output $@ $<

messageidentifier_string.go: messageidentifier.go tools/xtools/go.mod
	$(info generating $@...)
	@$(stringer) -type MessageIdentifier -trimprefix MessageIdentifier -output $@ $<

precision_string.go: precision.go tools/xtools/go.mod
	$(info generating $@...)
	@$(stringer) -type Precision -trimprefix Precision -output $@ $<

pkg/serial/baudrate_string.go: pkg/serial/baudrate.go tools/xtools/go.mod
	$(info generating $@...)
	@$(stringer) -type BaudRate -trimprefix BaudRate -output $@ $<

.PHONY: mockgen-generate
mockgen-generate: test/mocks/xsens/mocks.go

test/mocks/xsens/mocks.go: client.go go.mod
	go run github.com/golang/mock/mockgen \
		-destination $@ -package mockxsens \
		github.com/einride/xsens-go SerialPort
