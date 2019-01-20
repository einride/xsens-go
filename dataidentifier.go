package xsens

import "fmt"

// DataIdentifier is an Xsens data identifier.
//
// Each data identifier is constructed in this way:
//
//  +-------------------------------------------------------------------------------+
//  | 15 | 14 | 13 | 12 | 11 | 10 |  9 |  8 |  7 |  6 |  5 |  4 |  3 |  2 |  1 |  0 |
//  +------------------------+--------------+-------------------+-------------------+
//  | Group                  | Reserved     | Type              | Format            |
//  +------------------------+--------------+-------------------+-------------------+
//
// Group defines the category of the data, such as timestamps, orientations, angular velocities, etc.
//
// Type combined with Group defines the actual type of the data.
type DataIdentifier struct {
	DataType
	CoordinateSystem
	Precision
}

const (
	dataIdentifierTypeMask             uint16 = 0xf8f0
	dataIdentifierCoordinateSystemMask uint16 = 0x000c
	dataIdentifierPrecisionMask        uint16 = 0x0003
)

// MarshalUint16 returns the data identifier represented as a uint16.
func (d DataIdentifier) MarshalUint16() uint16 {
	return uint16(d.DataType) | uint16(d.CoordinateSystem) | uint16(d.Precision)
}

// MarshalUint16 sets the data identifier from a uint16 representation.
func (d *DataIdentifier) UnmarshalUint16(value uint16) {
	d.DataType = DataType(value & dataIdentifierTypeMask)
	d.CoordinateSystem = CoordinateSystem(value & dataIdentifierCoordinateSystemMask)
	d.Precision = Precision(value & dataIdentifierPrecisionMask)
}

// String returns a string representation of the data identifier.
func (d DataIdentifier) String() string {
	switch {
	case d.HasCoordinateSystem() && d.HasPrecision():
		return fmt.Sprintf("%v(%v,%v)", d.DataType, d.CoordinateSystem, d.Precision)
	case d.HasPrecision():
		return fmt.Sprintf("%v(%v)", d.DataType, d.Precision)
	default:
		return fmt.Sprintf("%v", d.DataType)
	}
}
