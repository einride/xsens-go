package serial

import (
	"time"
)

// Open a serial port.
func Open(filename string, baudRate BaudRate) (*Port, error) {
	impl, err := open(filename, baudRate)
	if err != nil {
		return nil, err
	}
	return &Port{impl: impl}, nil
}

// Port is a serial port.
type Port struct {
	impl *portImpl
}

// Flush the buffered data in the port.
func (p *Port) Flush() error {
	return p.impl.Flush()
}

// SetReadDeadline sets the port's read deadline.
func (p *Port) SetReadDeadline(t time.Time) error {
	return p.impl.SetReadDeadline(t)
}

// SetWriteDeadline sets the port's write deadline.
func (p *Port) SetWriteDeadline(t time.Time) error {
	return p.impl.SetWriteDeadline(t)
}

// Read data from the port.
func (p *Port) Read(b []byte) (n int, err error) {
	return p.impl.Read(b)
}

// Write data to the port.
func (p *Port) Write(b []byte) (int, error) {
	return p.impl.Write(b)
}

// Close the port.
func (p *Port) Close() error {
	return p.impl.Close()
}
