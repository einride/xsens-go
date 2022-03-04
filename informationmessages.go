package xsens

import (
	"encoding/binary"
	"fmt"
)

type DeviceID uint32

type ProductCode string

type HWVersion string // MAJOR.minor

func (d *DeviceID) UnmarshalBinary(data []byte) error {
	const mti100 = 4
	const mti600 = 8
	switch l := len(data); l {
	case mti100:
		*d = DeviceID(binary.BigEndian.Uint32(data))
	case mti600:
		*d = DeviceID(binary.BigEndian.Uint32(data[4:]))
	default:
		return fmt.Errorf("unexpected DeviceID length: want: (%d or %d), got: %d", mti100, mti600, l)
	}
	return nil
}

func (d *ProductCode) UnmarshalBinary(data []byte) error {
	*d = ProductCode(data)
	return nil
}

func (d *HWVersion) UnmarshalBinary(data []byte) error {
	if l := len(data); l != 2 {
		return fmt.Errorf("unexpected HWVersion length: want: %d, got: %d", 2, l)
	}
	*d = HWVersion(fmt.Sprintf("%d.%d", data[0], data[1]))
	return nil
}
