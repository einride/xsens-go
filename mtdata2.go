package xsens

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
)

type XDI uint16

const (
	packetCounter  XDI = 0x1020
	sampleTimeFine XDI = 0x1060
	quaternion     XDI = 0x2010
	ellipsoid      XDI = 0x5022
	latLon         XDI = 0x5042
	magnetic       XDI = 0xC000
	velocityXYZ    XDI = 0xd012
	statusWord     XDI = 0xe020
)

/*
should be integrated eventually
ErrorCodes = {
		0x03: "Invalid period",
		0x04: "Invalid message",
		0x1E: "Timer overflow",
		0x20: "Invalid baudrate",
		0x21: "Invalid parameter"
	}
*/

// XDIQuaternion contains orientation output expressed as a quaternion
type XDIQuaternion struct {
	Q0, Q1, Q2, Q3 float32
}

type XDILatLon struct {
	Lat, Lon float64
}

type XDIVelocityXYZ struct {
	VelX, VelY, VelZ float64
}

type XDIMagneticXYZ struct {
	MagX, MagY, MagZ float64
}

type xsens1632 struct {
	Fraction uint32
	Integer  int16
}

func (fp xsens1632) ToFloat() float64 {
	i := int64(fp.Integer) << 32
	i = i + int64(fp.Fraction)
	return float64(i) / math.Pow(2, 32)
}

type XsensData struct {
	PacketCouter               uint16
	SampleTimeFine, StatusWord uint32
	Altitude                   float64
	Quat                       XDIQuaternion
	Vel                        XDIVelocityXYZ
	Latlng                     XDILatLon
	Mag                        XDIMagneticXYZ
}

// MTData2Decode decode any number of packages from the given slice
func MTData2Decode(data []byte) (currentStatus XsensData, err error) {
	//var packetnum int
	buf := bytes.NewReader(data)

	for {
		var dataID XDI
		var dataLEN byte
		err = binary.Read(buf, binary.BigEndian, &dataID)
		if nil != err {
			if err == io.EOF {
				// no more data in pipe, controlled break
				err = nil
			} else {
				log.Printf("could not read packetcounter %+v", err)
			}
			return
		}
		err = binary.Read(buf, binary.BigEndian, &dataLEN)
		if nil != err {
			log.Printf("could not read data length %v", err)
			return
		}

		packetdata := make([]byte, dataLEN)
		err = binary.Read(buf, binary.BigEndian, &packetdata)
		if err != nil {
			log.Printf("could not get packetdata %+v", err)
			return
		}

		//fmt.Printf("Read MTData2 packet, id %#x, len %d, data %x\n",
		//	dataID, dataLEN, packetdata)

		packetBuf := bytes.NewReader(packetdata)

		switch dataID {
		case packetCounter:
			var packetCounter uint16
			err = binary.Read(packetBuf, binary.BigEndian, &packetCounter)
			if err != nil {
				log.Printf("could not read packetcounter %+v", err)
				return
			}
			currentStatus.PacketCouter = packetCounter
			break
		case sampleTimeFine:
			var sampleTimeFine uint32
			err = binary.Read(packetBuf, binary.BigEndian, &sampleTimeFine)
			if err != nil {
				log.Printf("could not read sampleTimeFine %+v", err)
				return
			}
			currentStatus.SampleTimeFine = sampleTimeFine
			break
		case quaternion:
			var quat XDIQuaternion
			err = binary.Read(packetBuf, binary.BigEndian, &quat)
			currentStatus.Quat = quat
			break
		case ellipsoid:
			var alt xsens1632
			err = binary.Read(packetBuf, binary.BigEndian, &alt)
			if err != nil {
				log.Printf("could not read ellipsoid altitude %+v", err)
				return
			}
			currentStatus.Altitude = alt.ToFloat()
			break
		case latLon:
			var lat, lon xsens1632
			err = binary.Read(packetBuf, binary.BigEndian, &lat)
			if err != nil {
				log.Printf("read lat error %+v", err)
				return
			}
			err = binary.Read(packetBuf, binary.BigEndian, &lon)
			if err != nil {
				log.Printf("read lat error %+v", err)
				return
			}
			currentStatus.Latlng = XDILatLon{lat.ToFloat(), lon.ToFloat()}
			break
		case velocityXYZ:
			var velx, vely, velz xsens1632
			err = binary.Read(packetBuf, binary.BigEndian, &velx)
			if err != nil {
				log.Printf("could not read velx %+v", err)
				return
			}
			err = binary.Read(packetBuf, binary.BigEndian, &vely)
			if err != nil {
				log.Printf("could not read vely %+v", err)
				return
			}
			err = binary.Read(packetBuf, binary.BigEndian, &velz)
			if err != nil {
				log.Printf("could not read velz %+v", err)
				return
			}
			currentStatus.Vel = XDIVelocityXYZ{velx.ToFloat(), vely.ToFloat(), velx.ToFloat()}
			break
		case statusWord:
			var statusWord uint32
			err = binary.Read(packetBuf, binary.BigEndian, &statusWord)
			if err != nil {
				log.Printf("could not read status word %+v", err)
				return
			}
			currentStatus.StatusWord = statusWord
			break
		case magnetic:
			var magx, magy, magz xsens1632
			err = binary.Read(packetBuf, binary.BigEndian, &magx)
			if err != nil {
				log.Printf("could not read x %+v", err)
				return
			}
			err = binary.Read(packetBuf, binary.BigEndian, &magy)
			if err != nil {
				log.Printf("could not read y %+v", err)
				return
			}
			err = binary.Read(packetBuf, binary.BigEndian, &magz)
			if err != nil {
				log.Printf("could not read z %+v", err)
				return
			}
			currentStatus.Mag = XDIMagneticXYZ{magx.ToFloat(), magy.ToFloat(), magz.ToFloat()}
			break
		default:
			fmt.Printf("\tUnknown packet ID %v\n", dataID)
		}
	}
	return
}
