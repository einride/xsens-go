package xsensgo

import (
	"log"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

var inputData = []byte{0x10, 0x10, 0x0C, 0x29, 0xFA, 0xC3, 0xE0, 0x07, 0xE2, 0x08, 0x09, 0x0D, 0x3B, 0x04, 0xF7, 0x10,
	0x20, 0x02, 0x43, 0x80, 0x20, 0x33, 0x18, 0x3F, 0xB6, 0x81, 0x2E, 0x60, 0x00, 0x00, 0x00, 0x40, 0x15, 0x3A, 0x60,
	0x80, 0x00, 0x00, 0x00, 0xC0, 0x06, 0x41, 0x85, 0xA0, 0x00, 0x00, 0x00, 0x40, 0x13, 0x18, 0xBF, 0x82, 0x88, 0x0F,
	0x20, 0x00, 0x00, 0x00, 0x3F, 0x28, 0xAD, 0xFA, 0x00, 0x00, 0x00, 0x00, 0x3F, 0xB9, 0x0D, 0xF2, 0x40, 0x00, 0x00,
	0x00, 0x40, 0x23, 0x18, 0xBF, 0xEC, 0xF4, 0x42, 0x80, 0x00, 0x00, 0x00, 0x3F, 0x93, 0x3E, 0x01, 0xC0, 0x00, 0x00,
	0x00, 0x40, 0x23, 0x92, 0xE5, 0xE0, 0x00, 0x00, 0x00, 0x40, 0x33, 0x18, 0x3F, 0x48, 0x9E, 0x3F, 0x00, 0x00, 0x00,
	0x00, 0x3F, 0x71, 0xFA, 0x9D, 0xC0, 0x00, 0x00, 0x00, 0x3F, 0x8B, 0xCA, 0xA0, 0x00, 0x00, 0x00, 0x00, 0x50, 0x23,
	0x08, 0x40, 0x59, 0x8F, 0x1D, 0x5D, 0x2D, 0x99, 0xF1, 0x50, 0x43, 0x10, 0x40, 0x4C, 0xD8, 0x6B, 0x0B, 0x41, 0x65,
	0x9C, 0x40, 0x27, 0xF1, 0xFE, 0x35, 0xEC, 0x57, 0xB8, 0x80, 0x23, 0x18, 0xBF, 0x40, 0xA0, 0x4B, 0xE0, 0x00, 0x00,
	0x01, 0xBF, 0x3B, 0x06, 0xBF, 0x00, 0x00, 0x00, 0x00, 0x3F, 0x55, 0x72, 0x9E, 0x80, 0x00, 0x00, 0x00, 0x80, 0x33,
	0x20, 0x3F, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xBE, 0xC5, 0x48, 0x0F, 0x00, 0x00, 0x00, 0x00, 0xBE, 0xC1,
	0x4B, 0xFF, 0x40, 0x00, 0x00, 0x00, 0x3E, 0xDB, 0x73, 0xFD, 0xE0, 0x00, 0x00, 0x00, 0xC0, 0x23, 0x18, 0x3F, 0xC8,
	0x6A, 0x38, 0xE0, 0x00, 0x00, 0x00, 0x3F, 0xD4, 0x09, 0xF6, 0x80, 0x00, 0x00, 0x00, 0xBF, 0xF3, 0x4D, 0x1D, 0xA0,
	0x00, 0x00, 0x00, 0xD0, 0x13, 0x18, 0x3F, 0xB2, 0x99, 0x03, 0xC0, 0x00, 0x00, 0x00, 0x3F, 0x87, 0x31, 0xDC, 0x40,
	0x00, 0x00, 0x00, 0x3F, 0x92, 0x86, 0x15, 0x80, 0x00, 0x00, 0x00, 0xE0, 0x20, 0x04, 0x01, 0x80, 0x00, 0x47}

var expectedData = Data{
	UTCTimestamp:    XDIUTCTime{NS: 704300000, Year: 2018, Month: 8, Day: 9, Hour: 13, Minute: 59, Second: 4, Conf: 247},
	PacketCounter:   17280,
	Euler:           XDIEulerAngles{Roll: 0.08790865, Pitch: 5.30700874, Yaw: -2.78199315},
	DeltaV:          XDIDeltaV{DVelX: -0.00904857, DVelY: 0.00018829, DVelZ: 0.09786905},
	Acc:             XDIAccelerationXYZ{AccX: -0.90481687, AccY: 0.01879122, AccZ: 9.78691006},
	FreeAcc:         XDIFreeAccelerationXYZ{FreeAccX: 0.00075129, FreeAccY: 0.00438940, FreeAccZ: 0.01357007},
	AltitudeMEllips: 102.23616723,
	Latlng:          XDILatLng{Lat: 57.69076672, Lng: 11.97264260},
	AngularVel:      XDIRateOfTurnXYZ{GyrX: -0.00050739, GyrY: -0.00041239, GyrZ: 0.00130907},
	DeltaQ:          XDIDeltaQ{DQ0: 1.00000000, DQ1: -0.00000254, DQ2: -0.00000206, DQ3: 0.00000655},
	Mag:             XDIMagneticXYZ{MagX: 0.19074164, MagY: 0.31310809, MagZ: -1.20632708},
	Vel:             XDIVelocityXYZ{VelX: 0.07264732, VelY: 0.01132557, VelZ: 0.01808961},
	StatusWord:      25165895,
}

func TestMTData2Decode(t *testing.T) {
	test := assert.New(t)

	currentStatus, err := Decode(inputData)
	test.Nil(err)

	log.Printf("currentstatus is %v, err is: %v", currentStatus, err)

	assert.Equal(t, expectedData.UTCTimestamp, currentStatus.UTCTimestamp)

	assert.Equal(t, expectedData.PacketCounter, currentStatus.PacketCounter)

	// Due to expected data being a float with fewer decimals, assert.equal cannot be used
	assert.True(t, math.Abs(expectedData.Euler.Roll-currentStatus.Euler.Roll) < 0.0000001)
	assert.True(t, math.Abs(expectedData.Euler.Pitch-currentStatus.Euler.Pitch) < 0.0000001)
	assert.True(t, math.Abs(expectedData.Euler.Yaw-currentStatus.Euler.Yaw) < 0.0000001)

	assert.True(t, math.Abs(expectedData.DeltaV.DVelX-currentStatus.DeltaV.DVelX) < 0.0000001)
	assert.True(t, math.Abs(expectedData.DeltaV.DVelY-currentStatus.DeltaV.DVelY) < 0.0000001)
	assert.True(t, math.Abs(expectedData.DeltaV.DVelZ-currentStatus.DeltaV.DVelZ) < 0.0000001)

	assert.True(t, math.Abs(expectedData.Acc.AccX-currentStatus.Acc.AccX) < 0.0000001)
	assert.True(t, math.Abs(expectedData.Acc.AccY-currentStatus.Acc.AccY) < 0.0000001)
	assert.True(t, math.Abs(expectedData.Acc.AccZ-currentStatus.Acc.AccZ) < 0.0000001)

	assert.True(t, math.Abs(expectedData.FreeAcc.FreeAccX-currentStatus.FreeAcc.FreeAccX) < 0.0000001)
	assert.True(t, math.Abs(expectedData.FreeAcc.FreeAccY-currentStatus.FreeAcc.FreeAccY) < 0.0000001)
	assert.True(t, math.Abs(expectedData.FreeAcc.FreeAccZ-currentStatus.FreeAcc.FreeAccZ) < 0.0000001)

	assert.True(t, math.Abs(expectedData.Vel.VelX-currentStatus.Vel.VelX) < 0.0000001)
	assert.True(t, math.Abs(expectedData.Vel.VelY-currentStatus.Vel.VelY) < 0.0000001)
	assert.True(t, math.Abs(expectedData.Vel.VelZ-currentStatus.Vel.VelZ) < 0.0000001)

	assert.True(t, math.Abs(expectedData.AngularVel.GyrX-currentStatus.AngularVel.GyrX) < 0.0000001)
	assert.True(t, math.Abs(expectedData.AngularVel.GyrY-currentStatus.AngularVel.GyrY) < 0.0000001)
	assert.True(t, math.Abs(expectedData.AngularVel.GyrZ-currentStatus.AngularVel.GyrZ) < 0.0000001)

	assert.True(t, math.Abs(expectedData.DeltaQ.DQ0-currentStatus.DeltaQ.DQ0) < 0.0000001)
	assert.True(t, math.Abs(expectedData.DeltaQ.DQ1-currentStatus.DeltaQ.DQ1) < 0.0000001)
	assert.True(t, math.Abs(expectedData.DeltaQ.DQ2-currentStatus.DeltaQ.DQ2) < 0.0000001)
	assert.True(t, math.Abs(expectedData.DeltaQ.DQ3-currentStatus.DeltaQ.DQ3) < 0.0000001)

	assert.True(t, math.Abs(expectedData.Mag.MagX-currentStatus.Mag.MagX) < 0.0000001)
	assert.True(t, math.Abs(expectedData.Mag.MagY-currentStatus.Mag.MagY) < 0.0000001)
	assert.True(t, math.Abs(expectedData.Mag.MagZ-currentStatus.Mag.MagZ) < 0.0000001)

	assert.Equal(t, expectedData.PacketCounter, currentStatus.PacketCounter)

	assert.True(t, math.Abs(expectedData.AltitudeMEllips-currentStatus.AltitudeMEllips) < 0.0000001)

	assert.True(t, math.Abs(expectedData.Latlng.Lat-currentStatus.Latlng.Lat) < 0.0000001)
	assert.True(t, math.Abs(expectedData.Latlng.Lng-currentStatus.Latlng.Lng) < 0.0000001)
}
