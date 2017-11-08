package xsens

import (
	"encoding/binary"
	"fmt"

	"io"

	"github.com/jacobsa/go-serial/serial"
	"log"
)

type header struct {
	Preamble byte
	BID      byte
	MID      byte
	LEN      byte
}

// All MIDs
const mtData2 = 54

var prt io.ReadWriteCloser

func readmsgs(callback func(data XsensData, err error)) (err error) {
	for {
		var data XsensData
		var h header
		err = binary.Read(prt, binary.BigEndian, &h)
		if nil != err {
			// might get io error here
			log.Println("Error reading from XSens: ", err)
			continue
		}

		if 0xfa != h.Preamble || 0xff != h.BID {
			err = fmt.Errorf("bad preamble and BID")
			callback(data, err)
			return
		}
		//log.Printf("Header: %x\n", h)
		var datalen uint16
		if 0xff > h.LEN {
			datalen = uint16(h.LEN)
		} else {
			err = binary.Read(prt, binary.BigEndian, &datalen)
			if nil != err {
				log.Println("Error reading datalength from MTMessage: ", err)
				callback(data, err)
				return
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
			// no more data, continue anyway
			continue
		}

		if nil != err {
			log.Println("Error reading data from XSens: ", err)
			callback(data, err)
			return
		}

		//log.Printf("Read %d bytes: %x\n", n, buf)

		// Read the checksum
		var checksum byte
		err = binary.Read(prt, binary.BigEndian, &checksum)
		if err != nil {
			log.Printf("could not read checksum %+v", err)
			return
		}
		// TODO: Validate chacksum

		if h.MID != mtData2 {
			err = fmt.Errorf("Unhandled MID %v\n", h.MID)
			log.Printf("%v", err)
			callback(data, err)
			return
		}

		data, err = mtData2Decode(buf)
		if err != nil {
			log.Printf("could not decode data %+v", err)
			callback(data, err)
			return
		}
		callback(data, nil)
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
		log.Printf("could not open ports %v", err)
		return
	}
	return
}

func Close() (err error) {
	if prt == nil {
		err = fmt.Errorf("could not close, no prt")
		log.Printf("%v", err)
		return
	}
	return prt.Close()
}

func Run(callback func(data XsensData, err error)) (err error) {
	err = Open()
	if err != nil {
		log.Printf("could not open %+v", err)
		return
	}
	defer Close()
	readmsgs(callback)
	return
}
