package xsensgo

import (
	"encoding/binary"
	"io"

	"os"

	"bytes"

	"log"

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

var xsensUSBHandle = func() string {
	handle := os.Getenv("XSENS_USB_HANDLE")
	if handle == "" {
		handle = "/dev/ttyUSB0"
	}
	return handle
}()

type Client struct {
	prt    io.ReadWriteCloser
	logger *zap.Logger
}

// Check "Client config" on drive for deepmap's config and more.
func (x *Client) readmsgs(callback func(data XsensData, err error)) error {
	for {
		// Read the header of the message
		h, err := readNextHeader(x.prt)
		if err != nil {
			return errors.Wrap(err, "could not read header")
		}

		var datalen uint16
		if h.LEN < 0xFF {
			datalen = uint16(h.LEN)
		} else {
			// If data package is of extended size. Will be this when following deepmap's setup + freeacc + mag.
			err = binary.Read(x.prt, binary.BigEndian, &datalen)
			if err != nil {
				err = errors.Wrap(err, "error reading datalength from MTMessage")
				callback(XsensData{}, err)
				return err
			}
		}

		// Create a buffer and read the whole data part into this buffer
		buf := make([]byte, datalen)
		var n int
		for n < int(datalen) && err == nil {
			var nn int
			nn, err = x.prt.Read(buf[n:])
			n += nn
		}

		if n >= int(datalen) {
			err = nil
			// no more data, continue anyway
			//log.Printf("no more data continue")
		}

		if err != nil {
			err = errors.Wrap(err, "error reading data from XSens")
			callback(XsensData{}, err)
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
			callback(XsensData{}, err)
			return err
		}

		// Check if message is GNSS
		if checkIfGNSS(buf) {
			// GNSS message, skip it
			continue
		}

		// Decode data in message
		data, err := mtData2Decode(buf)
		if err != nil {
			err = errors.Wrap(err, "could not decode data")
			callback(data, err)
			return err
		}

		// Set status from parsed data
		callback(data, nil)

	}
}

const (
	gnssID xdi = 0x7000
	group  xdi = 0xFF00
)

func checkIfGNSS(data []byte) bool {
	packets, err := parsePackets(bytes.NewReader(data))
	for i := 0; i < len(packets); i++ {
		if err != nil {
			log.Printf("Error parsing packets: %v", err)
			// TODO: Handle this error?
		}
		// Check if group ID is of type GNSS in any of packets
		if packets[i].id&group == gnssID {
			return true
		}
	}
	return false
}

func readNextHeader(prt io.Reader) (header, error) {
	h := header{}
	var data [3]byte

	// Wait for correct type of package
	for data[0] != busIdentifier {

		// Wait for the start of a package
		var sentry [1]byte
		for sentry[0] != packageStartIndicator {
			err := binary.Read(prt, binary.BigEndian, &sentry)
			if err != nil {
				return h, errors.Wrap(err, "could not read")
			}
		}

		err := binary.Read(prt, binary.BigEndian, &data)
		if err != nil {
			return h, errors.Wrap(err, "could not read")
		}
	}

	h.Preamble = packageStartIndicator
	h.BID = data[0]
	h.MID = data[1]
	h.LEN = data[2]

	return h, nil
}

func NewClient() (x *Client, err error) {
	// Configure and open the serial port to the Client
	options := serial.OpenOptions{
		PortName:        xsensUSBHandle,
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}
	prt, err := serial.Open(options)
	if err != nil {
		return nil, errors.Wrap(err, "could not open ports")
	}
	return &Client{prt: prt}, nil
}

func (x *Client) Close() (err error) {
	if x.prt == nil {
		err = errors.Errorf("could not close, no prt")
		return
	}
	return x.prt.Close()
}

func (x *Client) Run(callback func(data XsensData, err error)) (err error) {
	defer func() {
		err = x.Close()
		if err != nil {
			x.logger.Error("could not close")
		}
	}()
	return x.readmsgs(callback)
}

// Low-level message sending function.
func (x *Client) writeMsg(mid byte, data []byte) error {
	length := byte(len(data))
	var lendat []byte
	if length > 254 {
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

/*
## Send a message and read confirmation
	def write_ack(self, mid, data=[]):
		"""Send a message a read confirmation."""
		self.write_msg(mid, data)
		for tries in range(100):
			mid_ack, dataAck = self.read_msg()
			if MidAck==(mid+1):
				break
		else:
			raise MTException("Ack (0x%X) expected, MID 0x%X received instead"\
					" (after 100 tries)."%(mid+1, mid_ack))
		return dataAck
*/
type ack struct {
	MidAck, DataAck byte
}

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

func (x *Client) Reset() error {
	return x.writeAck(reset, []byte{})
}
