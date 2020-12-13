package xsens

// Precision is an Xsens data output precision.
type Precision uint8

//go:generate stringer -type Precision -trimprefix Precision

const (
	// PrecisionFloat32 uses single-precision IEEE 32-bit floating point numbers.
	PrecisionFloat32 Precision = 0x0

	// PrecisionFP1220 uses fixed point 12.20 32-bit numbers.
	PrecisionFP1220 Precision = 0x1

	// PrecisionFP1632 uses fixed point 16.32 48-bit numbers.
	PrecisionFP1632 Precision = 0x2

	// PrecisionFloat64 uses double-precision IEEE 64-bit floating point numbers.
	PrecisionFloat64 Precision = 0x3
)

// Size returns the size (in bytes) of a value with the current precision.
//
// Returns 0 for unsupported precisions.
func (p Precision) Size() uint8 {
	switch p {
	case PrecisionFloat32:
		return 4
	case PrecisionFP1220:
		return 4
	case PrecisionFP1632:
		return 6
	case PrecisionFloat64:
		return 8
	}
	return 0
}
