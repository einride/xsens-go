package xsens

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

const XDI_PacketCounter = 0x1020
const XDI_SampleTimeFine = 0x1060
const XDI_Quaternion = 0x2010
const XDI_AltitudeEllipsoid = 0x5022
const XDI_LatLon = 0x5042
const XDI_VelocityXYZ = 0xd012
const XDI_StatusWord = 0xe020
/*
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

type xsens1632 struct {
	Fraction uint32
	Integer  int16
}

func (fp xsens1632) ToFloat() float64 {
	i := int64(fp.Integer) << 32
	i = i + int64(fp.Fraction)
	return float64(i) / math.Pow(2, 32)
}

// MTData2Decode decode any number of packages from the given slice
func MTData2Decode(data []byte) (err error) {
	//var packetnum int
	buf := bytes.NewReader(data)
	for {
		var dataID uint16
		var dataLEN byte
		err = binary.Read(buf, binary.BigEndian, &dataID)
		if nil != err {
			return
		}
		err = binary.Read(buf, binary.BigEndian, &dataLEN)
		if nil != err {
			return
		}

		packetdata := make([]byte, dataLEN)
		err = binary.Read(buf, binary.BigEndian, &packetdata)

		fmt.Printf("Read MTData2 packet, id %#x, len %d, data %x\n",
			dataID, dataLEN, packetdata)

		packetBuf := bytes.NewReader(packetdata)

		switch dataID {
		case XDI_PacketCounter:
			var packetCounter uint16
			err = binary.Read(packetBuf, binary.BigEndian, &packetCounter)
			fmt.Printf("\tPacketCounter: %v\n", packetCounter)
		case XDI_SampleTimeFine:
			var sampleTimeFine uint32
			err = binary.Read(packetBuf, binary.BigEndian, &sampleTimeFine)
			fmt.Printf("\tSampleTimeFine: %v\n", sampleTimeFine)
		case XDI_Quaternion:
			var quat XDIQuaternion
			err = binary.Read(packetBuf, binary.BigEndian, &quat)
			fmt.Printf("\tQuaternion: %v\n", quat)
		case XDI_AltitudeEllipsoid:
			var alt xsens1632
			binary.Read(packetBuf, binary.BigEndian, &alt)
			altitude := alt.ToFloat()
			fmt.Printf("\tAltitude: %v\n", altitude)
		case XDI_LatLon:
			var lat, lon xsens1632
			binary.Read(packetBuf, binary.BigEndian, &lat)
			binary.Read(packetBuf, binary.BigEndian, &lon)
			latlon := XDILatLon{lat.ToFloat(), lon.ToFloat()}
			fmt.Printf("\tLatLon: %v\n", latlon)
		case XDI_VelocityXYZ:
			var velx, vely, velz xsens1632
			binary.Read(packetBuf, binary.BigEndian, &velx)
			binary.Read(packetBuf, binary.BigEndian, &vely)
			binary.Read(packetBuf, binary.BigEndian, &velz)
			vel := XDIVelocityXYZ{velx.ToFloat(), vely.ToFloat(), velx.ToFloat()}
			fmt.Printf("\tVelocityXYZ: %v\n", vel)
		case XDI_StatusWord:
			var statusWord uint32
			err = binary.Read(packetBuf, binary.BigEndian, &statusWord)
			fmt.Printf("\tStatusWord: %b\n", statusWord)
		default:
			fmt.Printf("\tUnknown packet ID %v\n", dataID)
		}
	}
}
