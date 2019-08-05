package xsens

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"golang.org/x/xerrors"
)

// field fixed lengths.
const (
	lengthOfPreamble          = 1
	lengthOfBusIdentifier     = 1
	lengthOfMessageIdentifier = 1
	lengthOfLength            = 1
	lengthOfExtendedLength    = 2
	lengthOfChecksum          = 1
)

// field max lengths.
const (
	minLengthOfExtendedData = 255
	maxLengthOfExtendedData = 2048
)

// MinLengthOfMessage is the minimum length of a valid Xsens message.
const MinLengthOfMessage = lengthOfPreamble +
	lengthOfBusIdentifier +
	lengthOfMessageIdentifier +
	lengthOfLength +
	lengthOfChecksum

// field start indices.
const (
	indexOfPreamble          = 0
	indexOfBusIdentifier     = indexOfPreamble + lengthOfPreamble
	indexOfMessageIdentifier = indexOfBusIdentifier + lengthOfBusIdentifier
	indexOfLength            = indexOfMessageIdentifier + lengthOfMessageIdentifier
	indexOfData              = indexOfLength + lengthOfLength
	indexOfExtendedLength    = indexOfLength + lengthOfLength
	indexOfExtendedData      = indexOfExtendedLength + lengthOfExtendedLength
)

// field constant values.
const (
	valueOfPreamble       = 0xfa
	valueOfBusIdentifier  = 0xff
	valueOfLengthExtended = 0xff
)

// Message represents a message sent to and from an Xsens device.
//
// A message has two basic structures; one with a standard length and one with extended
// length.
//
// Standard length
//
// The standard length message has a maximum of 254 data bytes and is used most frequently.
//
// An MT message (standard length) contains the following fields:
//
//  +----------+-----+-----+-----+------+----------+
//  | PREAMBLE | BID | MID | LEN | DATA | CHECKSUM |
//  +----------+-----+-----+-----+------+----------+
//  0          1     2     3     4
//
// Extended length
//
// In some cases the extended length message needs to be used if the number of data bytes exceeds
// 254 bytes.
//
// An MT message (extended length) contains these fields:
//
//  +----------+-----+-----+-----+---------+------+----------+
//  | PREAMBLE | BID | MID | LEN | EXT-LEN | DATA | CHECKSUM |
//  +----------+-----+-----+-----+---------+------+----------+
//  0          1     2     3     4         6
type Message []byte

// NewMessage creates a new Xsens message with the provided identifier and data.
//
// The provided data can be nil or empty, for messages without any data.
func NewMessage(mid MessageIdentifier, data []byte) Message {
	message := make(Message, MinLengthOfMessage+len(data))
	message[indexOfPreamble] = valueOfPreamble
	message[indexOfBusIdentifier] = valueOfBusIdentifier
	message[indexOfMessageIdentifier] = uint8(mid)
	message[indexOfLength] = uint8(len(data))
	copy(message[indexOfData:], data)
	message[len(message)-1] = 0xff & (-message.Checksum())
	return message
}

// String returns a string representation of the message.
func (m Message) String() string {
	if err := m.Validate(); err != nil {
		return fmt.Sprintf("InvalidMessage(%s)", hex.EncodeToString(m))
	}
	if m.IsError() {
		return fmt.Sprintf("%v(%v)", m.Identifier(), m.ErrorCode())
	}
	return fmt.Sprintf("%v(%s)", m.Identifier(), hex.EncodeToString(m.Data()))
}

// IsError is true if the message is an error message.
func (m Message) IsError() bool {
	return m.Identifier() == MessageIdentifierError && !m.IsExtended() && m.Length() == 1
}

// ErrorCode returns the message's error code.
//
// The error code is 0 (OK) for non-error messages.
func (m Message) ErrorCode() ErrorCode {
	if !m.IsError() {
		return ErrorCodeOK
	}
	return ErrorCode(m[indexOfData])
}

// Preamble of the message.
//
// Every message starts with the preamble. This field always contains the value 250 (=0xFA).
func (m Message) Preamble() uint8 {
	return m[indexOfPreamble]
}

// BusIdentifier of the message.
//
// An MT will only acknowledge a message (reply) if it is addressed with a valid BID. An MT will always
// acknowledge a message with the same BID that has been used to address it.
//
// Messages generated by the MT itself (i.e. not in acknowledge on a request) will always have a BID of 255 (0xFF).
// In practice, the only message for which this occurs is the MTData2 messages.
func (m Message) BusIdentifier() uint8 {
	return m[indexOfBusIdentifier]
}

// Identifier of the message.
func (m Message) Identifier() MessageIdentifier {
	return MessageIdentifier(m[indexOfMessageIdentifier])
}

// Length of the message.
func (m Message) Length() uint16 {
	if m.IsExtended() {
		return binary.BigEndian.Uint16(m[indexOfExtendedLength : indexOfExtendedLength+lengthOfExtendedLength])
	}
	return uint16(m[indexOfLength])
}

// IsExtended returns true if the message is extended.
func (m Message) IsExtended() bool {
	return m[indexOfLength] == valueOfLengthExtended
}

// Data returns the data of the message.
func (m Message) Data() []byte {
	if m.IsExtended() {
		return m[indexOfExtendedData : indexOfExtendedData+m.Length()]
	}
	return m[indexOfData : indexOfData+m.Length()]
}

// Checksum of the message.
//
// This field is used for communication error-detection. If all message bytes excluding the preamble are
// summed and the lower byte value of the result equals zero, the message is valid and it may be
// processed. The checksum value of the message should be included in the summation.
func (m Message) Checksum() uint8 {
	var checkSum uint8
	for i := indexOfBusIdentifier; i < len(m); i++ {
		checkSum += m[i]
	}
	return checkSum
}

// Validate the length and checksum of the message.
func (m Message) Validate() error {
	if len(m) < MinLengthOfMessage {
		return xerrors.Errorf("too few bytes: %v", len(m))
	}
	if m.Preamble() != valueOfPreamble {
		return xerrors.Errorf("invalid preamble: %v", m.Preamble())
	}
	if m.BusIdentifier() != valueOfBusIdentifier {
		return xerrors.Errorf("invalid bus identifier: %v", m.BusIdentifier())
	}
	if m.IsExtended() {
		if m.Length() < minLengthOfExtendedData || m.Length() > maxLengthOfExtendedData {
			return xerrors.Errorf("invalid extended length: %v", m.Length())
		}
	}
	if sum := m.Checksum(); sum != 0 {
		return xerrors.Errorf("invalid checksum: %v", sum)
	}
	return nil
}
