package xsens

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// OutputConfiguration is measurement data output configuration.
//
// The data is a list of maximum 32 data identifiers combined with a desired output frequency.
//
// For data that is sent with every MTData2 message (Timestamp, Status), the Output OutputFrequency will
// be ignored and will be set to 0xFFFF.
type OutputConfiguration []OutputConfigurationSetting

// Unmarshal sets *o from a wire representation of the output configuration.
func (o *OutputConfiguration) Unmarshal(data []byte) error {
	settingsCount := len(data) / 4
	if cap(*o) >= settingsCount {
		*o = (*o)[:settingsCount]
	} else {
		*o = append((*o)[:cap(*o)], make([]OutputConfigurationSetting, settingsCount-cap(*o))...)
	}
	for i := 0; i < settingsCount; i++ {
		(*o)[i].DataIdentifier.SetUint16(binary.BigEndian.Uint16(data[i*4:]))
		(*o)[i].OutputFrequency = OutputFrequency(binary.BigEndian.Uint16(data[i*4+2:]))
	}
	return nil
}

// Marshal returns the wire representation of the output configuration.
func (o *OutputConfiguration) Marshal() ([]byte, error) {
	buf := make([]byte, len(*o)*4)
	for i, setting := range *o {
		binary.BigEndian.PutUint16(buf[i*4:], setting.DataIdentifier.Uint16())
		binary.BigEndian.PutUint16(buf[i*4+2:], uint16(setting.OutputFrequency))
	}
	return buf, nil
}

// MarshalText returns a text representation of the output configuration.
func (o *OutputConfiguration) MarshalText() (string, error) {
	var buf bytes.Buffer
	for i, setting := range *o {
		if _, err := fmt.Fprintf(&buf, "%v:\n  %v", setting.DataIdentifier, setting.OutputFrequency); err != nil {
			return "", err
		}
		if i < len(*o)-1 {
			if err := buf.WriteByte('\n'); err != nil {
				return "", err
			}
		}
	}
	return buf.String(), nil
}

// OutputFrequency represents the output frequency of a specific Xsens measurement data type.
type OutputFrequency uint16

// MaxOutputFrequency is the sentinel value used for data types that should be included in every message, if possible.
const MaxOutputFrequency OutputFrequency = 0xffff

// String returns a string representation of the output frequency.
func (f OutputFrequency) String() string {
	switch f {
	case 0x0000, 0xffff:
		return "Max"
	default:
		return fmt.Sprintf("%d Hz", f)
	}
}

// OutputConfigurationSetting is the output configuration for a single measurement data type.
type OutputConfigurationSetting struct {
	// DataIdentifier is the data identifier of the data.
	DataIdentifier

	// OutputFrequency is the output frequency of the data.
	//
	// Selecting an Output OutputFrequency of either 0x0000 or 0xFFFF, makes the device select the
	// maximum frequency for the given data identifier. The device reports the resulting effective
	// frequencies in its response message.
	OutputFrequency
}
