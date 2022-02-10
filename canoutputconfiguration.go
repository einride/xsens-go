package xsens

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type CANOutputConfiguration []CANOutputConfigurationSetting

type CANIDLengthFlag bool

type IDMask = uint32

const (
	CANIDLengthFlag11bits CANIDLengthFlag = false
	CANIDLengthFlag29bits CANIDLengthFlag = true
)

const (
	canOutputCfgSettingSize = 8

	// CANDataIdentifier.
	canOutputCfgOffsetBlock1  = 0
	canOutputCfgLengthBlock1  = 1
	canOutputCfgCANDataIDMask = ^byte(1 << 7) // Masks reserved bit in canOutputCfg Block 1 MSB

	// CANIDLengthFlag.
	canOutputCfgOffsetBlock2     = canOutputCfgLengthBlock1
	canOutputCfgLengthBlock2     = 1
	canOutputCfgIDLengthFlagMask = byte(0b1) // Masks reserved bit in canOutputCfg Block 2 MSB

	// IDMask.
	canOutputCfgOffsetBlock3 = canOutputCfgLengthBlock2
	canOutputCfgLengthBlock3 = 4
	canOutputCfgIDMaskMask   = ^byte(0b111 << 5) // Masks reserved bit in canOutputCfg Block 3 MSB

	// OutputFrequency.
	canOutputCfgOffsetBlock4   = canOutputCfgLengthBlock3
	canOutputCfgLengthBlock4   = 2
	canOutputCfgOutputFreqMask = byte(0b111) // Masks reserved bit in canOutputCfg Block 4 MSB
)

// CANOutputConfigurationSetting is the output configuration for a single CAN measurement data type.
type CANOutputConfigurationSetting struct {
	// DataIdentifier is the data identifier of the data.
	CANDataIdentifier

	CANIDLengthFlag

	// TODO: Double check mask over interface
	IDMask IDMask

	// OutputFrequency is the output frequency of the data.
	//
	// Selecting an Output OutputFrequency of either 0x0000 or 0xFFFF, makes the device select the
	// maximum frequency for the given data identifier. The device reports the resulting effective
	// frequencies in its response message.
	OutputFrequency
}

// MarshalBinary returns the wire representation of the CAN output configuration.
func (o *CANOutputConfiguration) MarshalBinary() ([]byte, error) {
	buf := make([]byte, len(*o)*4)
	var b []byte
	// Push each setting to the buffer
	for i, setting := range *o {
		w := buf[i*canOutputCfgSettingSize : (i+1)*canOutputCfgSettingSize-1]

		// push CANDataIdentifier
		b = extractBytes(w, canOutputCfgOffsetBlock1, canOutputCfgLengthBlock1)
		b[0] = uint8(setting.CANDataIdentifier)
		b[0] &= canOutputCfgCANDataIDMask // mask MSB

		// push IDMaskLen
		b = extractBytes(w, canOutputCfgOffsetBlock2, canOutputCfgLengthBlock2)
		if setting.CANIDLengthFlag {
			b[0] = 1
		}

		// push IDMask
		b = extractBytes(w, canOutputCfgOffsetBlock3, canOutputCfgLengthBlock3)
		binary.BigEndian.PutUint32(b, setting.IDMask)
		b[0] &= canOutputCfgIDMaskMask // mask last byte

		// push OutputFrequency
		b = extractBytes(w, canOutputCfgOffsetBlock4, canOutputCfgLengthBlock4)
		binary.BigEndian.PutUint16(b, uint16(setting.OutputFrequency))
		b[0] &= canOutputCfgOutputFreqMask // mask last byte
	}
	return buf, nil
}

// MarshalText returns a text representation of the CAN output configuration.
func (o *CANOutputConfiguration) MarshalText() ([]byte, error) {
	var buf bytes.Buffer
	canIDLengthMap := map[CANIDLengthFlag]int{
		CANIDLengthFlag11bits: 11,
		CANIDLengthFlag29bits: 29,
	}
	for i, setting := range *o {
		if _, err := fmt.Fprintf(
			&buf,
			"%v:\n  %v Hz,\n  %v bits\n  %v mask",
			setting.CANDataIdentifier,
			setting.OutputFrequency,
			canIDLengthMap[setting.CANIDLengthFlag],
			setting.IDMask,
		); err != nil {
			return nil, err
		}
		if i < len(*o)-1 {
			if err := buf.WriteByte('\n'); err != nil {
				return nil, err
			}
		}
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary sets *o from a wire representation of the CAN output configuration.
func (o *CANOutputConfiguration) UnmarshalBinary(data []byte) error {
	settingsCount := len(data) / canOutputCfgSettingSize
	if cap(*o) >= settingsCount {
		*o = (*o)[:settingsCount]
	} else {
		*o = append((*o)[:cap(*o)], make([]CANOutputConfigurationSetting, settingsCount-cap(*o))...)
	}
	for i := 0; i < settingsCount; i++ {
		w := data[i*canOutputCfgSettingSize : (i+1)*canOutputCfgSettingSize-1]

		// get CANDataIdentifier
		b := copyBytes(w, canOutputCfgOffsetBlock1, canOutputCfgLengthBlock1)
		d := b[0] & canOutputCfgCANDataIDMask // mask MSB (single byte)
		(*o)[i].CANDataIdentifier = CANDataIdentifier(d)

		// get IDLengthFlag
		b = copyBytes(w, canOutputCfgOffsetBlock2, canOutputCfgLengthBlock2)
		d = b[0] & canOutputCfgIDLengthFlagMask // mask MSB (single byte)
		(*o)[i].CANIDLengthFlag = d != 0

		// get IDMask
		b = copyBytes(w, canOutputCfgOffsetBlock3, canOutputCfgLengthBlock3)
		b[0] &= canOutputCfgIDMaskMask // mask MSB
		(*o)[i].IDMask = binary.BigEndian.Uint32(b)

		// get OutputFrequency
		b = copyBytes(w, canOutputCfgOffsetBlock4, canOutputCfgLengthBlock4)
		b[0] &= canOutputCfgOutputFreqMask // mask MSB
		(*o)[i].OutputFrequency = OutputFrequency(binary.BigEndian.Uint16(b))
	}
	return nil
}

// IDMask returns the ID mask expected.
func (f CANIDLengthFlag) IDMask() IDMask {
	b := make([]byte, 0, 4)
	// bits 1-8
	b[0] = ^byte(0)

	// bits 9-11(or 16)
	if f == CANIDLengthFlag11bits {
		b[1] = 0b111
		return binary.LittleEndian.Uint32(b)
	}
	b[1] = ^byte(0)

	// bits 17-24
	b[2] = ^byte(0)

	// bits 25-29
	b[3] = 0b11111
	return binary.LittleEndian.Uint32(b)
}

// copyBytes returns a copy of the sub-slice of the byte array.
func copyBytes(w []byte, offset, length int) []byte {
	b := make([]byte, 0, length)
	copy(b, w[offset:offset+length])
	return b
}

// extractBytes returns a sub-slice of the byte array.
func extractBytes(w []byte, offset, length int) []byte {
	return w[offset : offset+length]
}
