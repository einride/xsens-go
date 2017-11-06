package xsens

import (
	"encoding/binary"
	"fmt"

	"io"

	"github.com/jacobsa/go-serial/serial"
	"log"
	"time"
	"math"
)

type header struct {
	Preamble byte
	BID      byte
	MID      byte
	LEN      byte
}

// All MIDs

const (
	timeout = 0.1
	sample = 100
)

var prt io.ReadWriteCloser

func readmsgs() (err error) {
	fmt.Println("Starting read loop...")
	for {
		var h header
		err = binary.Read(prt, binary.BigEndian, &h)
		if err != io.EOF {
			log.Printf("Error reading from XSens: %v", err)
		}

		if nil != err {
			log.Printf("Error reading from binary %v", err)
			return
		}

		if 0xfa == h.Preamble && 0xff == h.BID {
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

			switch MID(h.MID) {
			case MTData2:
				MTData2Decode(buf)
			default:
				fmt.Printf("Unhandled MID %v\n", h.MID)
			}
		}
	}
}

func flushStdIo() error {
	return nil
}

/*
## Place MT device in configuration mode.
	def GoToConfig(self):
		"""Place MT device in configuration mode."""
		self.write_ack(MID.GoToConfig)
*/
func goToConfig() error {
	_, err := writeAck(GoToConfig, []int{})
	return err
}

/*
## Send a message and read confirmation
	def write_ack(self, mid, data=[]):
		"""Send a message a read confirmation."""
		self.write_msg(mid, data)
		for tries in range(100):
			mid_ack, data_ack = self.read_msg()
			if mid_ack==(mid+1):
				break
		else:
			raise MTException("Ack (0x%X) expected, MID 0x%X received instead"\
					" (after 100 tries)."%(mid+1, mid_ack))
		return data_ack
*/
func writeAck(mid MID, data []int) (data_ack []byte, err error) {
	err = writeMsg(mid, data)
	if err != nil {
		log.Printf("could not write message %+v", err)
		return
	}

	for i := 0; i < 100; i++ {
		var mid_ack MID
		mid_ack, data_ack, err = readMsg()
		if err != nil {
			log.Printf("could not read msg %+v", err)
			return
		}
		if mid_ack == (mid + 1) {
			break
		}
		return
	}
	log.Printf("could not read data")
	err = fmt.Errorf("could not read msg after 100 tries")
	return
}

/*
def write_msg(self, mid, data=[]):
		"""Low-level message sending function."""
		length = len(data)
		if length>254:
			lendat = [0xFF, 0xFF&length, 0xFF&(length>>8)]
		else:
			lendat = [length]
		packet = [0xFA, 0xFF, mid] + lendat + list(data)
		packet.append(0xFF&(-(sum(packet[1:]))))
		msg = struct.pack('%dB'%len(packet), *packet)
		start = time.time()
		while (time.time()-start)<self.timeout and self.device.read():
			#print ".",
			pass
		self.device.write(msg)
		if verbose:
			print "MT: Write message id 0x%02X (%s) with %d data bytes: [%s]"%(mid, getMIDName(mid), length,
							' '.join("%02X"% v for v in data))
*/

func writeMsg(mid MID, data []int) error {
	length := len(data)
	lendat := []int{}
	if length > 254 {
		lendat = []int{0xFF, 0xFF&length, 0xFF&(length>>8)}
	} else {
		lendat = []int{length}
	}

	packet := append(append([]int{0xFA, 0xFF, int(mid)}, lendat...), data...)
	sum := 0
	for _, v := range packet[1:] {
		sum += v
	}
	packet = append(packet, 0xFF&(-(sum)))
	err := binary.Write(prt, binary.BigEndian, packet)
	if err != nil {
	    log.Printf("could not write to device %+v", err)
	    return err
	}
	return nil
}

func configureMTI(mtiSampleRate, mtiMode int) error {
	goToConfig()
	timeout := math.Pow(float64(mtiSampleRate), -1) + additionalTimeOutOffset
	mid := SetOutputConfiguration
	midReqDID := ReqDID
	dataReqDID := []int{0x00, 0x00}
	dataDID, err := writeAck(midReqDID, dataReqDID)
	if err != nil {
	    log.Printf("could not write ack %+v", err)
	    return err
	}
	masterID := binary.Write()

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
		log.Printf("could not open serial port %v", err)
		return
	}

	go readmsgs()
	return
}

func Close() (err error) {
	prt.Close()

	return
}

/*
def read_msg(self):
		"""Low-level message receiving function."""
		start = time.time()
		while (time.time()-start)<self.timeout:
			new_start = time.time()

			# Makes sure the buffer has 'size' bytes.
			def waitfor(size=1):
				while self.device.inWaiting() < size:
					if time.time()-new_start >= self.timeout:
						raise MTException("timeout waiting for message.")

			c = self.device.read()
			while (not c) and ((time.time()-new_start)<self.timeout):
				c = self.device.read()
			if not c:
				raise MTException("timeout waiting for message.")
			if ord(c)<>0xFA:
				continue
			# second part of preamble
			waitfor(3)
			if ord(self.device.read())<>0xFF:	# we assume no timeout anymore
				continue
			# read message id and length of message
			#msg = self.device.read(2)
			mid, length = struct.unpack('!BB', self.device.read(2))
			if length==255:	# extended length
				waitfor(2)
				length, = struct.unpack('!H', self.device.read(2))
			# read contents and checksum

			waitfor(length+1)
			buf = self.device.read(length+1)
			while (len(buf)<length+1) and ((time.time()-start)<self.timeout):
				buf+= self.device.read(length+1-len(buf))
			if (len(buf)<length+1):
				continue
			checksum = ord(buf[-1])
			data = struct.unpack('!%dB'%length, buf[:-1])
			if mid == MID.Error:
				sys.stderr.write("MT error 0x%02X: %s."%(data[0],
						MID.ErrorCodes[data[0]]))
			if verbose:
				print "MT: Got message id 0x%02X (%s) with %d data bytes: [%s]"%(mid, getMIDName(mid), length,
								' '.join("%02X"% v for v in data))
			if 0xFF&sum(data, 0xFF+mid+length+checksum):
				sys.stderr.write("invalid checksum; discarding data and "\
						"waiting for next message.\n")
				continue
			return (mid, buf[:-1])
		else:
			raise MTException("could not find message.")
*/

func readMsg() (m MID, data []byte, err error) {
	start := time.Now()
	for time.Now().Sub(start) < timeout {
		var h header
		err = binary.Read(prt, binary.BigEndian, &h)
		if err != nil {
			log.Printf("error reading from device %+v", err)
			return
		}

		buf := make([]byte, h.LEN)
		var n int
		for n < int(h.LEN) && err == nil {
			var nn int
			nn, err = prt.Read(buf[n:])
			n += nn
		}
		if n >= int(h.LEN) {
			err = nil
		} else if nil != err {
			fmt.Println("Error reading data from XSens: ", err)
			return
		}
		fmt.Printf("Read %d bytes: %x\n", n, buf)
		m = MID(h.MID)
		data = buf[:h.LEN]
		return
	}
	return
}

func setMtiOutputConfiguration(mtiSampleRate, mtiMode int) {
	flushStdIo()
	goToConfig()
	log.Printf("Device intiated at %d Hz", sample)
	configureMTI(mtiSampleRate, mtiMode)
}
