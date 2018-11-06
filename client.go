package xsensgo

import (
	"encoding/binary"
	"io"

	"github.com/jacobsa/go-serial/serial"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type header struct {
	Preamble byte
	BID      byte
	MID      byte
	LEN      byte
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

type ReceiverFunc func(data *Data, err error)

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

// Check "Client config" on drive for xsens's config and more.
func (x *Client) readMessages(callback ReceiverFunc) error {
	for {
		// Read the header of the message
		h, err := readNextHeader(x.prt)
		if err != nil {
			return errors.Wrap(err, "could not read header")
		}

		var dataLength uint16
		if h.LEN < 0xFF {
			dataLength = uint16(h.LEN)
		} else {
			// If data package is of extended size. Will be this when following deepmap's setup + freeacc + mag.
			err = binary.Read(x.prt, binary.BigEndian, &dataLength)
			if err != nil {
				err = errors.Wrap(err, "error reading datalength from MTMessage")
				callback(nil, err)
				return err
			}
		}

		// Create a buffer and read the whole data part into this buffer
		buf := make([]byte, dataLength)
		var n int
		for n < int(dataLength) && err == nil {
			var nn int
			nn, err = x.prt.Read(buf[n:])
			n += nn
		}

		if n >= int(dataLength) {
			err = nil
			// no more data, continue anyway
		}

		if err != nil {
			err = errors.Wrap(err, "error reading data from XSens")
			callback(nil, err)
			return err
		}

		// Read the checksum
		var checksum byte
		err = binary.Read(x.prt, binary.BigEndian, &checksum)
		if err != nil {
			return errors.Wrap(err, "could not read checksum")
		}
		// TODO: Validate chacksum

		// Check if Message ID is of type mtData2
		if h.MID != mtData2 {
			err = errors.Errorf("Unhandled MID %v\n", h.MID)
			callback(nil, err)
			return err
		}

		// Check if message is GNSS
		if checkIfGNSS(buf) {
			// GNSS message, skip it
			continue
		}

		// Decode data in message
		data, err := Decode(buf)
		if err != nil {
			err = errors.Wrap(err, "could not decode data")
			callback(data, err)
			return err
		}

		// Set status from parsed data
		callback(data, nil)
	}
}

func readNextHeader(prt io.Reader) (*header, error) {
	var data [3]byte

	// Wait for correct type of package
	for data[0] != busIdentifier {

		// Wait for the start of a package
		var sentry [1]byte
		for sentry[0] != packageStartIndicator {
			err := binary.Read(prt, binary.BigEndian, &sentry)
			if err != nil {
				return nil, errors.Wrap(err, "could not read")
			}
		}

		err := binary.Read(prt, binary.BigEndian, &data)
		if err != nil {
			return nil, errors.Wrap(err, "could not read")
		}
	}

	h := &header{}
	h.Preamble = packageStartIndicator
	h.BID = data[0]
	h.MID = data[1]
	h.LEN = data[2]

	return h, nil
}

// Close the serial port connection
func (x *Client) Close() (err error) {
	if x.prt == nil {
		return errors.Errorf("could not close, no prt")
	}
	return x.prt.Close()
}

func (x *Client) Run(callback ReceiverFunc) (err error) {
	defer func() {
		err = x.Close()
		if err != nil {
			x.logger.Error("could not close")
		}
	}()
	return x.readMessages(callback)
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
