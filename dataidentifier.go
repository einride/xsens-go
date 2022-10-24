package xsens

import (
	"fmt"
)

// DataIdentifier is an Xsens data identifier.
//
// Each data identifier is constructed in this way:
//
// +-------------------------------------------------------------------------------+
// | 15 | 14 | 13 | 12 | 11 | 10 |  9 |  8 |  7 |  6 |  5 |  4 |  3 |  2 |  1 |  0 |
// +------------------------+--------------+-------------------+-------------------+
// | Group                  | Reserved     | Type              | Format            |
// +------------------------+--------------+-------------------+-------------------+
//
// Group defines the category of the data, such as timestamps, orientations, angular velocities, etc.
//
// Type combined with Group defines the actual type of the data.
type DataIdentifier struct {
	DataType         DataType
	CoordinateSystem CoordinateSystem
	Precision        Precision
}

const (
	dataIdentifierTypeMask             uint16 = 0xf8f0
	dataIdentifierCoordinateSystemMask uint16 = 0x000c
	dataIdentifierPrecisionMask        uint16 = 0x0003
)

// Uint16 returns the data identifier represented as a uint16.
func (d DataIdentifier) Uint16() uint16 {
	return uint16(d.DataType) | uint16(d.CoordinateSystem) | uint16(d.Precision)
}

// SetUint16 sets the data identifier from a uint16 representation.
func (d *DataIdentifier) SetUint16(value uint16) {
	d.DataType = DataType(value & dataIdentifierTypeMask)
	d.CoordinateSystem = CoordinateSystem(value & dataIdentifierCoordinateSystemMask)
	d.Precision = Precision(value & dataIdentifierPrecisionMask)
}

// DataSize returns the data size (in bytes) of measurement data with the current identifier.
//
// Returns 0 for unsupported data identifiers.
func (d DataIdentifier) DataSize() uint8 {
	switch d.DataType {
	// scalars
	case DataTypeTemperature,
		DataTypeAltitudeEllipsoid:
		return d.Precision.Size()
	// vectors
	case DataTypeDeltaV,
		DataTypeAcceleration,
		DataTypeFreeAcceleration,
		DataTypeAccelerationHR,
		DataTypeRateOfTurn,
		DataTypeRateOfTurnHR,
		DataTypeEulerAngles,
		DataTypePositionECEF,
		DataTypeVelocityXYZ,
		DataTypeMagneticField:
		return 3 * d.Precision.Size()
	case DataTypeDeltaQ, DataTypeQuaternion:
		return 4 * d.Precision.Size()
	case DataTypeUTCTime:
		return 12
	case DataTypePacketCounter:
		return 2
	case DataTypeSampleTimeFine:
		return 4
	case DataTypeSampleTimeCoarse:
		return 4
	case DataTypeRotationMatrix:
		return 9 * d.Precision.Size()
	case DataTypeBaroPressure:
		return 4
	case DataTypeLatLon:
		return 2
	case DataTypeGNSSPVTData:
		return 76
	case DataTypeGNSSSatInfo:
		return 8 // plus variable number of satellites
	case DataTypeStatusByte:
		return 1
	case DataTypeStatusWord:
		return 4
	}
	return 0
}

// String returns a string representation of the data identifier.
func (d DataIdentifier) String() string {
	switch {
	case d.DataType.HasCoordinateSystem() && d.DataType.HasPrecision():
		return fmt.Sprintf("%v(%v,%v)", d.DataType, d.CoordinateSystem, d.Precision)
	case d.DataType.HasPrecision():
		return fmt.Sprintf("%v(%v)", d.DataType, d.Precision)
	default:
		return fmt.Sprintf("%v", d.DataType)
	}
}
