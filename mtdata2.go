package xsens

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	bp "github.com/roman-kachanovsky/go-binary-pack/binary-pack"
	"strings"
	"log"
)

type XDI uint16

const (
	packetCounter     XDI = 0x1020
	sampleTimeFine    XDI = 0x1060
	quaternion        XDI = 0x2010
	altitudeEllipsoid XDI = 0x5022
	latLon            XDI = 0x5042
	velocityXYZ       XDI = 0xd012
	statusWord        XDI = 0xe020
	magneticField     XDI = 0xC020
)

type MID uint16

const (
	MTData2    MID = 54
	Error      MID = 0x42
	GoToConfig MID = 0x30
	SetOutputConfiguration MID = 0xC0
	ReqDID MID = 0x00
)

const additionalTimeOutOffset = 0.010

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
		var dataID XDI
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
		case packetCounter:
			var packetCounter uint16
			err = binary.Read(packetBuf, binary.BigEndian, &packetCounter)
			fmt.Printf("\tPacketCounter: %v\n", packetCounter)
		case sampleTimeFine:
			var sampleTimeFine uint32
			err = binary.Read(packetBuf, binary.BigEndian, &sampleTimeFine)
			fmt.Printf("\tSampleTimeFine: %v\n", sampleTimeFine)
		case quaternion:
			var quat XDIQuaternion
			err = binary.Read(packetBuf, binary.BigEndian, &quat)
			fmt.Printf("\tQuaternion: %v\n", quat)
		case altitudeEllipsoid:
			var alt xsens1632
			binary.Read(packetBuf, binary.BigEndian, &alt)
			altitude := alt.ToFloat()
			fmt.Printf("\tAltitude: %v\n", altitude)
		case latLon:
			var lat, lon xsens1632
			binary.Read(packetBuf, binary.BigEndian, &lat)
			binary.Read(packetBuf, binary.BigEndian, &lon)
			latlon := XDILatLon{lat.ToFloat(), lon.ToFloat()}
			fmt.Printf("\tLatLon: %v\n", latlon)
		case velocityXYZ:
			var velx, vely, velz xsens1632
			binary.Read(packetBuf, binary.BigEndian, &velx)
			binary.Read(packetBuf, binary.BigEndian, &vely)
			binary.Read(packetBuf, binary.BigEndian, &velz)
			vel := XDIVelocityXYZ{velx.ToFloat(), vely.ToFloat(), velx.ToFloat()}
			fmt.Printf("\tVelocityXYZ: %v\n", vel)
		case statusWord:
			var statusWord uint32
			err = binary.Read(packetBuf, binary.BigEndian, &statusWord)
			fmt.Printf("\tStatusWord: %b\n", statusWord)
		case magneticField:
			err = binary.Read(packetBuf, binary.BigEndian, &statusWord)
		default:
			fmt.Printf("\tUnknown packet ID %v\n", dataID)
		}
	}
}

/*
def parse_GNSS(data_id, content, ffmt):
			o = {}
			pvtFlag = False
			if (data_id&0x00F0) == 0x10:	# GNSS PVT DATA
				o['iTOW'],x1,x2,x3,x4,x5,x6,x7,x8,x9,o['fix'],o['flag'],o['nSat'],x10,lon,lat,h,a,hAcc, \
				vAcc,vN,vE,vD,x11,x12,sAcc,headAcc,headVeh,gDop,pDop,tDop,vDop,hDop,nDop,eDop = \
						struct.unpack('!LHBBBBBBLiBBBBiiiiLLiiiiiLLIHHHHHHH', content)
				o['lat'], o['lon'], o['hEll'], o['hMsl'], o['velN'], o['velE'], o['velD'], \
				o['horzAcc'], o['vertAcc'], o['speedAcc'], o['GDOP'],  o['PDOP'],  o['TDOP'],\
				o['VDOP'], o['HDOP'], o['NDOP'], o['EDOP'], o['heading'], o['headingAcc'] = 1e-7*lat, 1e-7*lon, 1e-3*h, \
						1e-3*a, 1e-3*vN, 1e-3*vE, 1e-3*vD, 1e-3*hAcc, 1e-3*vAcc, 1e-3*sAcc, 1e-2*gDop, \
						1e-2*pDop, 1e-2*tDop, 1e-2*vDop, 1e-2*hDop, 1e-2*nDop, 1e-2*eDop, 1e-5*headVeh, 1e-5*headAcc
				pvtFlag = True
			elif (data_id&0x00F0) == 0x20:	# GNSS SAT Info
				o['iTOW'], o['numCh'] = struct.unpack('!LBxxx', content[:8])
				channels = []
				ch = {}
				for i in range(o['numCh']):
					ch['gnssId'], ch['svId'], ch['cno'], ch['flags'] = \
							struct.unpack('!BBBB', content[8+4*i:12+4*i])
					channels.append(ch)
					ch = {} # empty
				o['channels'] = channels
			else:
				raise MTException("unknown packet: 0x%04X."%data_id)
			return o, pvtFlag
*/
func parseGNSS(data_id XDI, content []byte, ffmt string) error {
	o := map[string]string{}
	pvtFlag := false
	if (data_id & 0x00F0) == 0x10 { // GNSS PVT DATA
		unPack := bp.BinaryPack{}
		format := strings.Split("!LHBBBBBBLiBBBBiiiiLLiiiiiLLIHHHHHHH", "")
		res, err := unPack.UnPack(format, content)
		if err != nil {
		    log.Printf(" %+v", err)
		    return err
		}
		o["iTOW"],x1,x2,x3,x4,x5,x6,x7,x8,x9,o["fix"],o["flag"],o["nSat"],x10,lon,lat,h,a,hAcc, vAcc,vN,vE,vD,x11,x12,sAcc,headAcc,headVeh,gDop,pDop,tDop,vDop,hDop,nDop,eDop := res...

		o["lat"], o["lon"], o["hEll"], o["hMsl"], o["velN"], o["velE"], o["velD"],
			o["horzAcc"], o["vertAcc"], o["speedAcc"], o["GDOP"],  o["PDOP"],  o["TDOP"],
			o["VDOP"], o["HDOP"], o["NDOP"], o["EDOP"], o["heading"], o["headingAcc"] = 1e-7*lat, 1e-7*lon, 1e-3*h,
			1e-3*a, 1e-3*vN, 1e-3*vE, 1e-3*vD, 1e-3*hAcc, 1e-3*vAcc, 1e-3*sAcc, 1e-2*gDop,
			1e-2*pDop, 1e-2*tDop, 1e-2*vDop, 1e-2*hDop, 1e-2*nDop, 1e-2*eDop, 1e-5*headVeh, 1e-5*headAcc
			pvtFlag = true
	}
	return nil
}
