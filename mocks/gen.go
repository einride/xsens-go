package test

//go:generate go run -mod=mod github.com/golang/mock/mockgen -package mockserial -destination mockserial/mocks.go go.bug.st/serial Port
