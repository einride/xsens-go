// +build linux
// +build go1.12

package serial

import (
	"os"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
	"golang.org/x/xerrors"
)

func open(filename string, baudRate BaudRate) (*portImpl, error) {
	fd, err := syscall.Open(filename, unix.O_RDWR|unix.O_NOCTTY|unix.O_NONBLOCK, 0666)
	if err != nil {
		return nil, xerrors.Errorf("open serial port %s: %w", filename, err)
	}
	// Important: Setting non-block in Go >=1.12 registers the fd with the Go runtime poller.
	//            This makes deadlines and wake-up on close work.
	if err := syscall.SetNonblock(fd, true); err != nil {
		return nil, err
	}
	termiosBaudRate, err := toTermios(baudRate)
	if err != nil {
		return nil, xerrors.Errorf("open serial port %s: %w", filename, err)
	}
	t := unix.Termios{
		// From: https://www.cmrr.umn.edu/~strupp/serial.html
		//
		// The cflag member contains two options that should always be enabled, CLOCAL and CREAD.
		// These will ensure that your program does not become the 'owner' of the port subject to sporatic job control
		// and hangup signals, and also that the serial interface driver will read incoming data bytes.
		//
		// From: https://en.wikipedia.org/wiki/Serial_port#Data_bits
		//
		// 8 data bits are almost universally used in newer applications,
		Cflag:  uint32(unix.CREAD) | uint32(unix.CLOCAL) | uint32(unix.CS8),
		Ispeed: termiosBaudRate,
		Ospeed: termiosBaudRate,
	}
	// From man termios(3)
	//
	// MIN > 0, TIME == 0 (blocking read)
	//   read(2) blocks until MIN bytes are available, and returns up to the number of bytes requested.
	t.Cc[unix.VMIN] = 1
	t.Cc[unix.VTIME] = 0
	// From: https://www.cmrr.umn.edu/~strupp/serial.html
	//
	// TCSETS: Sets the serial port settings immediately.
	if err := unix.IoctlSetTermios(fd, unix.TCSETS, &t); err != nil {
		return nil, xerrors.Errorf("open serial port %s: %w", filename, err)
	}
	f := os.NewFile(uintptr(fd), filename)
	c, err := f.SyscallConn()
	if err != nil {
		return nil, xerrors.Errorf("open serial port %s: %w", filename, err)
	}
	return &portImpl{f: f, c: c}, nil
}

type portImpl struct {
	f *os.File
	c syscall.RawConn
}

func (p *portImpl) Flush() error {
	var errFlush error
	err := p.c.Control(func(fd uintptr) {
		errFlush = unix.IoctlSetInt(int(fd), unix.TCFLSH, unix.TCIOFLUSH)
	})
	if err != nil {
		return xerrors.Errorf("serial port %s: flush: %w", p.f.Name(), err)
	}
	if errFlush != nil {
		return xerrors.Errorf("serial port %s: flush: %w", p.f.Name(), err)
	}
	return nil
}

func (p *portImpl) SetReadDeadline(t time.Time) error {
	if err := p.f.SetReadDeadline(t); err != nil {
		return xerrors.Errorf("serial port %s: set read deadline: %w", p.f.Name(), err)
	}
	return nil
}

func (p *portImpl) SetWriteDeadline(t time.Time) error {
	if err := p.f.SetWriteDeadline(t); err != nil {
		return xerrors.Errorf("serial port %s: set write deadline: %w", p.f.Name(), err)
	}
	return nil
}

func (p *portImpl) Read(b []byte) (int, error) {
	n, err := p.f.Read(b)
	if err != nil {
		return n, xerrors.Errorf("serial port %s: read: %w", p.f.Name(), err)
	}
	return n, nil
}

func (p *portImpl) Write(b []byte) (int, error) {
	n, err := p.f.Write(b)
	if err != nil {
		return n, xerrors.Errorf("serial port %s: write: %w", p.f.Name(), err)
	}
	return n, nil
}

func (p *portImpl) Close() error {
	if err := p.f.Close(); err != nil {
		return xerrors.Errorf("serial port %s: close: %w", p.f.Name(), err)
	}
	return nil
}

func toTermios(b BaudRate) (uint32, error) {
	switch b {
	case BaudRate4800:
		return unix.B4800, nil
	case BaudRate9600:
		return unix.B9600, nil
	case BaudRate19200:
		return unix.B19200, nil
	case BaudRate38400:
		return unix.B38400, nil
	case BaudRate57600:
		return unix.B57600, nil
	case BaudRate115200:
		return unix.B115200, nil
	case BaudRate230400:
		return unix.B230400, nil
	case BaudRate460800:
		return unix.B460800, nil
	case BaudRate921600:
		return unix.B921600, nil
	case BaudRate2000000:
		return unix.B2000000, nil
	default:
		return 0, xerrors.Errorf("baud rate to termios: unsupported value: %v", b)
	}
}
