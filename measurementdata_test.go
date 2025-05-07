package xsens

import (
	"fmt"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestConvert_Scalar(t *testing.T) {
	for _, tt := range []DataIdentifier{
		{
			DataType:  DataTypeTemperature,
			Precision: PrecisionFloat32,
		},
		{
			DataType:  DataTypeTemperature,
			Precision: PrecisionFP1220,
		},
		{
			DataType:  DataTypeTemperature,
			Precision: PrecisionFP1632,
		},
		{
			DataType:  DataTypeTemperature,
			Precision: PrecisionFloat64,
		},
	} {
		t.Run(fmt.Sprintf("precision %v", tt), func(t *testing.T) {
			var scalar Scalar = 1
			data, err := scalar.MarshalMTData2Packet(tt)
			assert.NilError(t, err)
			var newScalar Scalar
			err = newScalar.UnmarshalMTData2Packet(data)
			assert.Equal(t, err, nil)
			assert.Equal(t, scalar, newScalar)
		})
	}
}

func TestConvert_VectorXYZ(t *testing.T) {
	for _, tt := range []DataIdentifier{
		{
			DataType:  DataTypeVelocityXYZ,
			Precision: PrecisionFloat32,
		},
		{
			DataType:  DataTypeVelocityXYZ,
			Precision: PrecisionFP1220,
		},
		{
			DataType:  DataTypeVelocityXYZ,
			Precision: PrecisionFP1632,
		},
		{
			DataType:  DataTypeVelocityXYZ,
			Precision: PrecisionFloat64,
		},
	} {
		t.Run(fmt.Sprintf("vector %v", tt), func(t *testing.T) {
			vec := VectorXYZ{X: 1, Y: 2, Z: 3}
			data, err := vec.MarshalMTData2Packet(tt)
			assert.NilError(t, err)
			var newVec VectorXYZ
			err = newVec.UnmarshalMTData2Packet(data)
			assert.NilError(t, err)
			assert.Equal(t, vec, newVec)
		})
	}
}

func TestConvert_Quaternion(t *testing.T) {
	for _, tt := range []DataIdentifier{
		{
			DataType:  DataTypeQuaternion,
			Precision: PrecisionFloat32,
		},
		{
			DataType:  DataTypeQuaternion,
			Precision: PrecisionFP1220,
		},
		{
			DataType:  DataTypeQuaternion,
			Precision: PrecisionFP1632,
		},
		{
			DataType:  DataTypeQuaternion,
			Precision: PrecisionFloat64,
		},
	} {
		t.Run(fmt.Sprintf("quat %v", tt), func(t *testing.T) {
			vec := Quaternion{Q0: 1, Q1: 2, Q2: 3, Q3: 4}
			data, err := vec.MarshalMTData2Packet(tt)
			assert.NilError(t, err)
			var n Quaternion
			err = n.UnmarshalMTData2Packet(data)
			assert.NilError(t, err)
			assert.Equal(t, vec, n)
		})
	}
}

func TestConvert_RotationMatrix(t *testing.T) {
	for _, tt := range []DataIdentifier{
		{
			DataType:  DataTypeRotationMatrix,
			Precision: PrecisionFloat32,
		},
		{
			DataType:  DataTypeRotationMatrix,
			Precision: PrecisionFP1220,
		},
		{
			DataType:  DataTypeRotationMatrix,
			Precision: PrecisionFP1632,
		},
		{
			DataType:  DataTypeRotationMatrix,
			Precision: PrecisionFloat64,
		},
	} {
		t.Run(fmt.Sprintf("rotation matrix %v", tt), func(t *testing.T) {
			org := RotationMatrix{
				A: 1,
				B: 2,
				C: 3,
				D: 4,
				E: 5,
				F: 6,
				G: 7,
				H: 8,
				I: 9,
			}
			data, err := org.MarshalMTData2Packet(tt)
			assert.NilError(t, err)
			var n RotationMatrix
			err = n.UnmarshalMTData2Packet(data)
			assert.NilError(t, err)
			assert.Equal(t, org, n)
		})
	}
}

func TestConvert_LatLon(t *testing.T) {
	for _, tt := range []DataIdentifier{
		{
			DataType:  DataTypeLatLon,
			Precision: PrecisionFloat32,
		},
		{
			DataType:  DataTypeLatLon,
			Precision: PrecisionFP1220,
		},
		{
			DataType:  DataTypeLatLon,
			Precision: PrecisionFP1632,
		},
		{
			DataType:  DataTypeLatLon,
			Precision: PrecisionFloat64,
		},
	} {
		t.Run(fmt.Sprintf("rotation matrix %v", tt), func(t *testing.T) {
			org := LatLon{
				Lat: 1,
				Lon: 2,
			}
			data, err := org.MarshalMTData2Packet(tt)
			assert.NilError(t, err)
			var n LatLon
			err = n.UnmarshalMTData2Packet(data)
			assert.NilError(t, err)
			assert.Equal(t, org, n)
		})
	}
}

func TestConvert_StatusByte(t *testing.T) {
	const dataType = DataTypeStatusByte
	for _, tt := range []DataIdentifier{
		{
			DataType:  dataType,
			Precision: PrecisionFloat32,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFP1220,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFP1632,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFloat64,
		},
	} {
		t.Run(fmt.Sprintf("rotation matrix %v", tt), func(t *testing.T) {
			org := StatusByte(1)
			data, err := org.MarshalMTData2Packet(tt)
			assert.NilError(t, err)
			var n StatusByte
			err = n.UnmarshalMTData2Packet(data)
			assert.NilError(t, err)
			assert.Equal(t, org, n)
		})
	}
}

func TestConvert_GNSSPVTData(t *testing.T) {
	const dataType = DataTypeGNSSPVTData
	for _, tt := range []DataIdentifier{
		{
			DataType:  dataType,
			Precision: PrecisionFloat32,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFP1220,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFP1632,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFloat64,
		},
	} {
		t.Run(fmt.Sprintf("GNSSPVTData %v", tt), func(t *testing.T) {
			org := GNSSPVTData{
				ITOW:      1,
				Year:      3,
				Month:     4,
				Day:       5,
				Hour:      6,
				Min:       7,
				Sec:       8,
				Valid:     UTCDateValidFlag | UTCTimeOfDayValidFlag | UTCTimeOfDayValidFlag,
				TAcc:      10,
				Nano:      11,
				FixType:   12,
				Flags:     13,
				NumSV:     14,
				Reserved1: 15,
				Lon:       17,
				Lat:       18,
				Height:    19,
				HMSL:      20,
				HAcc:      21,
				VAcc:      22,
				VelN:      23,
				VelE:      24,
				VelD:      25,
				GSpeed:    26,
				HeadMot:   27,
				SAcc:      28,
				HeadAcc:   29,
				HeadVeh:   30,
				GDOP:      31,
				PDOP:      32,
				TDOP:      33,
				VDOP:      34,
				HDOP:      35,
				NDOP:      36,
				EDOP:      37,
			}
			data, err := org.MarshalMTData2Packet(tt)
			assert.NilError(t, err)
			var n GNSSPVTData
			err = n.UnmarshalMTData2Packet(data)
			assert.NilError(t, err)
			assert.Equal(t, org, n)
		})
	}
}

func TestConvert_SampleTimeFine(t *testing.T) {
	const dataType = DataTypeSampleTimeFine
	for _, tt := range []DataIdentifier{
		{
			DataType:  dataType,
			Precision: PrecisionFloat32,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFP1220,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFP1632,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFloat64,
		},
	} {
		t.Run(fmt.Sprintf("SampleTimeFine %v", tt), func(t *testing.T) {
			org := SampleTimeFine(1)
			data, err := org.MarshalMTData2Packet(tt)
			assert.NilError(t, err)
			var n SampleTimeFine
			err = n.UnmarshalMTData2Packet(data)
			assert.NilError(t, err)
			assert.Equal(t, org, n)
		})
	}
}

func TestConvert_SampleTimeCoarse(t *testing.T) {
	const dataType = DataTypeSampleTimeCoarse
	for _, tt := range []DataIdentifier{
		{
			DataType:  dataType,
			Precision: PrecisionFloat32,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFP1220,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFP1632,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFloat64,
		},
	} {
		t.Run(fmt.Sprintf("SampleTimeCoarse %v", tt), func(t *testing.T) {
			org := SampleTimeCoarse(1)
			data, err := org.MarshalMTData2Packet(tt)
			assert.NilError(t, err)
			var n SampleTimeCoarse
			err = n.UnmarshalMTData2Packet(data)
			assert.NilError(t, err)
			assert.Equal(t, org, n)
		})
	}
}

func TestClient_UTCTime(t *testing.T) {
	ts := time.Date(1, 2, 3, 4, 5, 6, 7, time.UTC)
	u := UTCTime{}
	u.UnmarshalTime(ts)
	assert.Equal(t, ts, u.Time())
}

func TestConvert_UTCTime(t *testing.T) {
	const dataType = DataTypeUTCTime
	for _, tt := range []DataIdentifier{
		{
			DataType:  dataType,
			Precision: PrecisionFloat32,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFP1220,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFP1632,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFloat64,
		},
	} {
		t.Run(fmt.Sprintf("UTCTime %v", tt), func(t *testing.T) {
			org := UTCTime{
				Ns:     1,
				Year:   2,
				Month:  3,
				Day:    4,
				Hour:   5,
				Minute: 6,
				Second: 7,
				Valid:  UTCDateValidFlag | UTCTimeOfDayValidFlag | UTCTimeOfDayFullyResolvedFlag,
			}
			data, err := org.MarshalMTData2Packet(tt)
			assert.NilError(t, err)
			var n UTCTime
			err = n.UnmarshalMTData2Packet(data)
			assert.NilError(t, err)
			assert.Equal(t, org, n)
		})
	}
}

func TestConvert_PacketCounter(t *testing.T) {
	const dataType = DataTypePacketCounter
	for _, tt := range []DataIdentifier{
		{
			DataType:  dataType,
			Precision: PrecisionFloat32,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFP1220,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFP1632,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFloat64,
		},
	} {
		t.Run(fmt.Sprintf("PacketCounter %v", tt), func(t *testing.T) {
			org := PacketCounter(43875)
			data, err := org.MarshalMTData2Packet(tt)
			assert.NilError(t, err)
			var n PacketCounter
			err = n.UnmarshalMTData2Packet(data)
			assert.NilError(t, err)
			assert.Equal(t, org, n)
		})
	}
}

func TestConvert_GNSSSatInfo(t *testing.T) {
	const dataType = DataTypeGNSSSatInfo
	for _, tt := range []DataIdentifier{
		{
			DataType:  dataType,
			Precision: PrecisionFloat32,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFP1220,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFP1632,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFloat64,
		},
	} {
		t.Run(fmt.Sprintf("GNSSSatInfo %v", tt), func(t *testing.T) {
			org := GNSSSatInfo{
				ITOW:   1,
				NumSVS: 2,
				Res1:   3,
				Res2:   4,
				Res3:   5,
			}
			data, err := org.MarshalMTData2Packet(tt)
			assert.NilError(t, err)
			var n GNSSSatInfo
			err = n.UnmarshalMTData2Packet(data)
			assert.NilError(t, err)
			assert.Equal(t, org, n)
		})
	}
}

func TestConvert_StatusWord(t *testing.T) {
	const dataType = DataTypeStatusWord
	for _, tt := range []DataIdentifier{
		{
			DataType:  dataType,
			Precision: PrecisionFloat32,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFP1220,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFP1632,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFloat64,
		},
	} {
		t.Run(fmt.Sprintf("StatusWord %v", tt), func(t *testing.T) {
			org := StatusWord(1)
			data, err := org.MarshalMTData2Packet(tt)
			assert.NilError(t, err)
			var n StatusWord
			err = n.UnmarshalMTData2Packet(data)
			assert.NilError(t, err)
			assert.Equal(t, org, n)
		})
	}
}

func TestConvert_BaroPressure(t *testing.T) {
	const dataType = DataTypeBaroPressure
	for _, tt := range []DataIdentifier{
		{
			DataType:  dataType,
			Precision: PrecisionFloat32,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFP1220,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFP1632,
		},
		{
			DataType:  dataType,
			Precision: PrecisionFloat64,
		},
	} {
		t.Run(fmt.Sprintf("BaroPressure %v", tt), func(t *testing.T) {
			org := BaroPressure(1)
			data, err := org.MarshalMTData2Packet(tt)
			assert.NilError(t, err)
			var n BaroPressure
			err = n.UnmarshalMTData2Packet(data)
			assert.NilError(t, err)
			assert.Equal(t, org, n)
		})
	}
}
