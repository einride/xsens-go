package xsensemulator

import (
	"fmt"
	"io"
	"net"
	"time"
)

type UDPSerialPort struct {
	io.ReadWriteCloser
	opts            *udpSerialPortOptions
	OriginConn      *net.UDPConn
	DestinationAddr *net.UDPAddr
}

func NewUDPSerialPort(
	origin string,
	destination string,
	udpSerialOpts ...UDPSerialPortOption,
) (*UDPSerialPort, error) {
	opts := defaultOptions()
	for _, udpSerialOpt := range udpSerialOpts {
		udpSerialOpt(opts)
	}

	udpOriginAddr, err := net.ResolveUDPAddr("udp", origin)
	if err != nil {
		return nil, fmt.Errorf("new udp serial port: %w", err)
	}
	udpDestinationAddr, err := net.ResolveUDPAddr("udp", destination)
	if err != nil {
		return nil, fmt.Errorf("new udp serial port: %w", err)
	}
	originConn, err := net.ListenUDP("udp", udpOriginAddr)
	if err != nil {
		return nil, fmt.Errorf("new udp serial port: %w", err)
	}
	return &UDPSerialPort{
		opts:            opts,
		OriginConn:      originConn,
		DestinationAddr: udpDestinationAddr,
	}, nil
}

func (t UDPSerialPort) Read(p []byte) (n int, err error) {
	// Check if a timeout opts have been added
	if t.opts.timeout != time.Duration(0) {
		err = t.OriginConn.SetReadDeadline(time.Now().Add(t.opts.timeout))
		if err != nil {
			return 0, fmt.Errorf("udp serial port read: %w", err)
		}
	}
	n, _, err = t.OriginConn.ReadFromUDP(p)
	return n, err
}

func (t UDPSerialPort) Write(p []byte) (n int, err error) {
	// Check if a timeout opts have been added
	if t.opts.timeout != time.Duration(0) {
		err = t.OriginConn.SetWriteDeadline(time.Now().Add(t.opts.timeout))
		if err != nil {
			return 0, fmt.Errorf("udp serial port write: %w", err)
		}
	}
	return t.OriginConn.WriteToUDP(p, t.DestinationAddr)
}

func (t UDPSerialPort) Close() error {
	return t.OriginConn.Close()
}

type udpSerialPortOptions struct {
	// Timeout for setting read/write deadlines
	timeout time.Duration
}

// defaultTransmitterOptions returns transmitter options with sensible default values.
func defaultOptions() *udpSerialPortOptions {
	return &udpSerialPortOptions{}
}

// UDPSerialPortOption configures an UDPSerialPort.
type UDPSerialPortOption func(*udpSerialPortOptions)

// WithTransmitInterface configures the interface to transmit on.
func WithTimeout(timeout time.Duration) UDPSerialPortOption {
	return func(opt *udpSerialPortOptions) {
		opt.timeout = timeout
	}
}
