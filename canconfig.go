package xsens

import (
	"fmt"
)

type CANConfig struct {
	Enable   bool
	BaudRate CANBaudRate
}

const (
	canCfgEnableOffset   = 2
	canCfgBaudrateOffset = 3
	canCfgEnableMask     = byte(0b1)     // Masks reserved bit in Enable byte
	canCfgBaudrateMask   = ^byte(1 << 7) // Masks reserved bit in BaudRate byte
)

// MarshalBinary returns the wire representation of the CAN configuration.
func (o *CANConfig) MarshalBinary() ([]byte, error) {
	result := make([]byte, 4)
	if o.Enable {
		result[canCfgEnableOffset] = 1
	}
	result[canCfgBaudrateOffset] = uint8(o.BaudRate) & canCfgBaudrateMask
	return result, nil
}

// MarshalText returns a text representation of the CAN configuration.
func (o *CANConfig) MarshalText() ([]byte, error) {
	s := fmt.Sprintf("Enable: %v, BaudRate: %v\n", o.Enable, o.BaudRate)
	return []byte(s), nil
}

// UnmarshalBinary sets *o from a wire representation of the CAN configuration.
func (o *CANConfig) UnmarshalBinary(data []byte) error {
	if o == nil {
		return fmt.Errorf("cannot unmarshal to a nil pointer")
	}
	o.Enable = (data[canCfgEnableOffset] & canCfgEnableMask) == 1
	o.BaudRate = CANBaudRate(data[canCfgBaudrateOffset] & canCfgBaudrateMask)
	return nil
}
