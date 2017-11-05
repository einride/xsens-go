package xsens

import (
	"encoding/binary"
	"fmt"

	"io"

	"github.com/jacobsa/go-serial/serial"
)

type header struct {
	Preamble byte
	BID      byte
	MID      byte
	LEN      byte
}

// All MIDs
const MTData2 = 54

var prt io.ReadWriteCloser

func readmsgs() (err error) {
	fmt.Println("Starting read loop...")
	for {
		var h header
		err = binary.Read(prt, binary.BigEndian, &h)
		if nil != err {
			if err != io.EOF {
				fmt.Println("Error reading from XSens: ", err)
			}
		} else if 0xfa == h.Preamble && 0xff == h.BID {
			fmt.Printf("Header: %x\n", h)
			var datalen uint16
			if 0xff > h.LEN {
				datalen = uint16(h.LEN)
			} else {
				err = binary.Read(prt, binary.BigEndian, &datalen)
				if nil != err {
					fmt.Println("Error reading datalength from MTMessage: ", err)
				}
			}

			// Create a buffer and read the whole data part into this buffer
			buf := make([]byte, datalen)
			var n int
			for n < int(datalen) && err == nil {
				var nn int
				nn, err = prt.Read(buf[n:])
				n += nn
			}
			if n >= int(datalen) {
				err = nil
			} else if nil != err {
				fmt.Println("Error reading data from XSens: ", err)
			}
			fmt.Printf("Read %d bytes: %x\n", n, buf)

			// Read the checksum
			var checksum byte
			err = binary.Read(prt, binary.BigEndian, &checksum)

			// TODO: Validate chacksum

			switch h.MID {
			case MTData2:
				MTData2Decode(buf)
			default:
				fmt.Printf("Unhandled MID %v\n", h.MID)
			}
		}
	}
}

func Open() (err error) {
	// Configure and open the serial port to the xsens
	options := serial.OpenOptions{
		PortName:        "/dev/ttyUSB0",
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}
	prt, err = serial.Open(options)
	if nil != err {
		return
	}

	go readmsgs()
	return
}

func Close() (err error) {
	prt.Close()

	return
}
