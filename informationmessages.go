package xsens

import (
	"encoding/binary"
	"fmt"
)

type DeviceID uint32

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
