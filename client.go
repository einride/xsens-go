package xsens

import (
	"encoding/binary"
	"io"

	"github.com/jacobsa/go-serial/serial"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type header struct {
	BID byte
	MID byte
	LEN byte
}

func (h *header) Read(r io.Reader) error {
	var data [3]byte
	// Wait for correct type of package
	for data[0] != busIdentifier {
		// Wait for the start of a package
		var sentry [1]byte
		for sentry[0] != packageStartIndicator {
			err := binary.Read(r, binary.BigEndian, &sentry)
			if err != nil {
				return errors.Wrap(err, "could not read")
			}
		}
		err := binary.Read(r, binary.BigEndian, &data)
		if err != nil {
			return errors.Wrap(err, "could not read")
		}
	}
	h.BID = data[0]
	h.MID = data[1]
	h.LEN = data[2]
	return nil
}

// All MIDs
const (
	mtData2               = 54
	reset                 = 0x40
	packageStartIndicator = 0xFA
	busIdentifier         = 0xFF
)

type Client struct {
	prt    io.ReadWriteCloser
	logger *zap.Logger
}

// create a new client, handle is usually /dev/ttyUSB0 on linux systems
func NewClient(prt io.ReadWriteCloser, logger *zap.Logger) (x *Client) {
	// Configure and open the serial port to the Client
	return &Client{prt: prt, logger: logger}
}

func DefaultSerialPort() (io.ReadWriteCloser, error) {
	return serial.Open(serial.OpenOptions{
		PortName:        "/dev/ttyUSB0",
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	})
}

// Close the serial port connection
func (x *Client) Close() (err error) {
	if x.prt == nil {
		return errors.Errorf("could not close, no prt")
	}
	return x.prt.Close()
}

// Low-level message sending function.
func (x *Client) writeMsg(mid byte, data []byte) error {
	length := byte(len(data))
	var lendat []byte
	if length > 0xFE {
		//lendat = []byte{
		//	0xFF,
		//	0xFF & length,
		//	0xFF & (length >> 8), // gives warning length (8 bits) too small for shift of 8
		//}$
		return errors.Errorf("this probably never happens, and currently should not be used")
	}
	lendat = []byte{length}
	packet := append(append([]byte{0xFA, 0xFF, mid}, lendat...), data...)
	sum := byte(0)
	for _, s := range packet[1:] {
		sum += s
	}
	packet = append(packet, 0xFF&(-sum))
	err := binary.Write(x.prt, binary.BigEndian, &packet)
	return err
}

type ack struct {
	MidAck, DataAck byte
}

// Send a message and read confirmation
func (x *Client) writeAck(mid byte, data []byte) error {
	err := x.writeMsg(mid, data)
	if err != nil {
		return err
	}

	for i := 1; i <= 100; i++ {
		a := ack{}
		err = binary.Read(x.prt, binary.BigEndian, &a)
		if err != nil {
			return errors.Wrap(err, "could not read")
		}
		if a.MidAck == mid+1 {
			return nil
		}
	}
	return errors.Errorf("no ack got")
}

// Write a rest to the xsens
func (x *Client) Reset() error {
	return x.writeAck(reset, []byte{})
}
