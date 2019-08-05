// +build !linux !go1.12

package serial

import (
	"runtime"
	"time"
)

func notSupported() string {
	return "not supported on " + runtime.GOOS + " and " + runtime.Version()
}

func open(_ string, _ BaudRate) (*portImpl, error) {
	panic(notSupported())
}

type portImpl struct{}

func (portImpl) Flush() error {
	panic(notSupported())
}

func (portImpl) SetReadDeadline(_ time.Time) error {
	panic(notSupported())
}

func (portImpl) SetWriteDeadline(_ time.Time) error {
	panic(notSupported())
}

func (portImpl) Read(_ []byte) (int, error) {
	panic(notSupported())
}

func (portImpl) Write(_ []byte) (int, error) {
	panic(notSupported())
}

func (portImpl) Close() error {
	panic(notSupported())
}
