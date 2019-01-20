package xsens

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

const (
	packetDataIdentifierLength = 2
	packetSizeLength           = 1
	packetDataIdentifierStart  = 0
	packetDataLengthStart      = packetDataIdentifierStart + packetDataIdentifierLength
	packetDataStart            = packetDataLengthStart + packetSizeLength
)

// MTData2 represents an Xsens MTData2 payload.
type MTData2 []byte

func (m MTData2) String() string {
	return fmt.Sprintf("MTData2(%s)", hex.EncodeToString(m))
}

// PacketAt returns the MTData2 packet starting at the provided index.
func (m MTData2) PacketAt(i int) (MTData2Packet, error) {
	if len(m) < i+packetDataIdentifierLength+packetSizeLength {
		return nil, errors.New("insufficient data")
	}
	packetDataLength := int(m[i+packetDataLengthStart])
	if len(m) < i+packetDataIdentifierLength+packetSizeLength+packetDataLength {
		return nil, errors.New("insufficient data")
	}
	return MTData2Packet(m[i : i+packetDataIdentifierLength+packetSizeLength+packetDataLength]), nil
}

// MTData2Packet represents an individual packet of an XSens MTData2 message.
type MTData2Packet []byte

// String returns a string representation of the packet.
func (m MTData2Packet) String() string {
	return fmt.Sprintf("MTData2Packet(%s)", hex.EncodeToString(m))
}

// Identifier returns the packet's data identifier.
func (m MTData2Packet) Identifier() DataIdentifier {
	var identifier DataIdentifier
	identifier.UnmarshalUint16(
		binary.BigEndian.Uint16(m[packetDataIdentifierStart : packetDataIdentifierStart+packetDataIdentifierLength]),
	)
	return identifier
}

// Data returns the packet data.
func (m MTData2Packet) Data() []byte {
	return m[packetDataStart:]
}
