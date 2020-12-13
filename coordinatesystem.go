package xsens

// CoordinateSystem represents the coordinate system of a measurement data output.
type CoordinateSystem uint8

//go:generate stringer -type CoordinateSystem -trimprefix CoordinateSystem

// Coordinate systems.
const (
	CoordinateSystemEastNorthUp   CoordinateSystem = 0x0
	CoordinateSystemNorthEastDown CoordinateSystem = 0x4
	CoordinateSystemNorthWestUp   CoordinateSystem = 0x8
)
