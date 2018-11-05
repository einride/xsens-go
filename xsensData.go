package xsensgo

import (
	"bytes"
	"log"
	"math"
)

// group
type xdi uint16

const (
	timestamp       xdi = 0x1000
	orientationData xdi = 0x2000
	acceleration    xdi = 0x4000
	position        xdi = 0x5000
	gnss            xdi = 0x7000
	angularVelocity xdi = 0x8000
	magnetic        xdi = 0xC000
	velocity        xdi = 0xD000
	statusWord      xdi = 0xE000
	gnssID          xdi = 0x7000
	group           xdi = 0xFF00
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

type Data struct {
	PacketCounter   uint16
	StatusWord      uint32
	UTCTimestamp    XDIUTCTime
	AltitudeMEllips float64
	Euler           XDIEulerAngles
	Vel             XDIVelocityXYZ
	Latlng          XDILatLng
	Acc             XDIAccelerationXYZ
	FreeAcc         XDIFreeAccelerationXYZ
	Mag             XDIMagneticXYZ
	AngularVel      XDIRateOfTurnXYZ
	DeltaQ          XDIDeltaQ
	DeltaV          XDIDeltaV
	Quat            XDIQuaternion
}

type XDIDeltaV struct {
	DVelX, DVelY, DVelZ float64
}

type XDIQuaternion struct {
	Q0, Q1, Q2, Q3 float64
}

// Contains the delta quaternion value of the SDI output.
type XDIDeltaQ struct {
	DQ0, DQ1, DQ2, DQ3 float64
}

type XDIEulerAngles struct {
	Roll, Pitch, Yaw float64
}

type XDILatLng struct {
	Lat, Lng float64
}

type XDIVelocityXYZ struct {
	VelX, VelY, VelZ float64
}

type XDIAccelerationXYZ struct {
	AccX, AccY, AccZ float64
}
type XDIFreeAccelerationXYZ struct {
	FreeAccX, FreeAccY, FreeAccZ float64
}

type XDIMagneticXYZ struct {
	MagX, MagY, MagZ float64
}

type XDIRateOfTurnXYZ struct {
	GyrX, GyrY, GyrZ float64
}

type XDIUTCTime struct {
	NS                                     int32
	Year                                   uint16
	Month, Day, Hour, Minute, Second, Conf uint8
}

/*
heading=arctan(Yh/Xh);
if    (Xh<0)         {heading=180-heading;}
elseif(Xh>0  && Yh<0){heading=-heading;}
elseif(Xh>0  && Yh>0){heading=360-heading}
elseif(Xh==0 && Yh<0){heading=90}
elseif(Xh==0 && Yh>0){heading=270}*/

func (data Data) Heading() (heading float64) {
	cRoll := math.Cos(data.Euler.Roll * math.Pi / 180)
	sRoll := math.Sin(data.Euler.Roll * math.Pi / 180)
	cPitch := math.Cos(data.Euler.Pitch * math.Pi / 180)
	sPitch := math.Sin(data.Euler.Pitch * math.Pi / 180)

	mX := data.Mag.MagX
	mY := data.Mag.MagY
	mZ := data.Mag.MagZ

	// Tilt compensated magnetometer values
	magX := mX*cPitch + mZ*sPitch
	magY := mY*sRoll*sPitch + mY*cRoll - mZ*sRoll*cPitch

	heading = math.Atan2(-magY, magX)*180/math.Pi + 90
	// Adjust for magnetic declination (Gothenburg)
	heading += 3.71
	if heading > 180.0 {
		heading -= 360.0
	} else if heading < -180.0 {
		heading += 360.0
	}

	return
}

/*
velocity = (VelY-VelX)x(VelZ-VelX) = - VelX - VelY + VelZ
velocity output from mt in m/s, Velocity() output in km/h
*/
// Absolute velocity
func (vel XDIVelocityXYZ) Velocity() float64 {
	return math.Abs(vel.VelX + vel.VelY + vel.VelZ)
}

/*
acceleration = (AccY-AccX)x(AccZ-AccX) = - AccX - AccY + AccZ
Acceleration output from mt in m/s2, Acceleration() output in km/h2
*/

// Calculating absolute acceleration, with gravity bias
func (acc XDIAccelerationXYZ) Acceleration() float64 {
	return math.Abs(acc.AccX + acc.AccY + acc.AccZ)
}

// Calculating absolute acceleration, without gravity bias
func (freeAcc XDIFreeAccelerationXYZ) Acceleration() float64 {
	return math.Abs(freeAcc.FreeAccX + freeAcc.FreeAccY + freeAcc.FreeAccZ)
}

// Make Euler to Quaternion conversion according to JPL convention
func (euler XDIEulerAngles) ToQuaternion() (quaternion XDIQuaternion) {
	cYaw := math.Cos(euler.Yaw * 0.5)
	sYaw := math.Sin(euler.Yaw * 0.5)
	cRoll := math.Cos(euler.Roll * 0.5)
	sRoll := math.Sin(euler.Roll * 0.5)
	cPitch := math.Cos(euler.Pitch * 0.5)
	sPitch := math.Cos(euler.Pitch * 0.5)
	qW := cYaw*cRoll*cPitch + sYaw*sRoll*sPitch
	qX := cYaw*sRoll*cPitch - sYaw*cRoll*sPitch
	qY := cYaw*cRoll*sPitch + sYaw*sRoll*cPitch
	qZ := sYaw*cRoll*cPitch - cYaw*sRoll*sPitch
	return XDIQuaternion{Q0: qW, Q1: qX, Q2: qY, Q3: qZ}
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

func checkIfGNSS(data []byte) bool {
	packets, err := parsePackets(bytes.NewReader(data))
	if err != nil {
		log.Printf("Error parsing packets: %v", err)
		// TODO: Handle this error?
	}

	for i := 0; i < len(packets); i++ {
		// Check if group ID is of type GNSS in any of packets
		if packets[i].id&group == gnssID {
			return true
		}
	}
	return false
}
