package xsens

// FixType represents an Xsens GNSS fix type.
type FixType uint8

const (
	FixTypeNoFix                FixType = 0x00
	FixTypeDeadReckoningOnly    FixType = 0x01
	FixType2DFix                FixType = 0x02
	FixType3DFix                FixType = 0x03
	FixTypeGNSSAndDeadReckoning FixType = 0x04
	FixTypeTimeOnly             FixType = 0x05
)
