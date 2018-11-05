package xsensgo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
)

type packet struct {
	id   xdi
	data []byte
}

// MTData2Decode decode any number of packages from the given slice.
func Decode(data []byte) (*Data, error) {
	// Split data into list of packets to be decoded
	packets, err := parsePackets(bytes.NewReader(data))

	if err != nil {
		return nil, err
	}
	// Decode content of packets
	return decodePackets(packets)
}

func parsePackets(reader io.Reader) ([]packet, error) {
	var err error
	var packets []packet
	for {
		// Read data id of packet
		var id xdi
		e := binary.Read(reader, binary.BigEndian, &id)
		// If there's no more data in pipe, we don't actually have an error, we just want to stop reading
		if e != nil {
			if e == io.EOF {
				break
			} else {
				// Note that an unexpected EOF here is an actual error
				return packets, errors.Wrap(e, "could not read id")
			}
		}

		var dataLEN byte
		// Read length of packet
		e = binary.Read(reader, binary.BigEndian, &dataLEN)
		if e != nil {
			if !isAnyKindOfEOF(e) {
				return packets, errors.Wrap(e, fmt.Sprintf("could not read data length for packet id='%v'", id))
			}
		}

		packetData := make([]byte, dataLEN)
		// Read contents of packet
		e = binary.Read(reader, binary.BigEndian, &packetData)
		if e != nil {
			if !isAnyKindOfEOF(e) {
				return packets, errors.Wrap(e, fmt.Sprintf("could not get packetData for packet id='%v'", id))
			}
		}
		packets = append(packets, packet{id: id, data: packetData})
	}

	return packets, err
}

func isAnyKindOfEOF(e error) bool {
	return strings.Contains(e.Error(), io.EOF.Error())
}

func decodePackets(packets []packet) (*Data, error) {
	var currentStatus Data

	// Decode all packets separately
	for _, packet := range packets {
		group := packet.id & 0xFF00         // Group is defined as first byte of the packet id
		groupTypeName := packet.id & 0x00F0 // Group type is defined as bit 8-11 of the packet id

		//logger.Info("group: %v, type: %v, data length %v, packet data len: %v", group, groupTypeName, len(data), dataLEN)
		var err error
		switch group {
		case timestamp:
			err = decodeTimestamp(groupTypeName, packet, &currentStatus)
		case orientationData:
			err = decodeOrientation(groupTypeName, packet, &currentStatus)
		case position:
			err = decodePosition(groupTypeName, packet, &currentStatus)
		case velocity:
			err = decodeVelocity(groupTypeName, packet, &currentStatus)
		case statusWord:
			err = decodeStatusWord(groupTypeName, packet, &currentStatus)
		case acceleration:
			err = decodeAcceleration(groupTypeName, packet, &currentStatus)
		case magnetic:
			err = decodeMagnetic(groupTypeName, packet, &currentStatus)
		case gnss:
			switch groupTypeName {
			case 0x10: // GNSS PVT Data
				// Do nothing as this currently is not interesting
				break
			case 0x20: // GNSS Satellites info
				// Do nothing as this currently is not interesting
				break
			}
		case angularVelocity:
			err = decodeAngularVelocity(groupTypeName, packet, &currentStatus)
		default:
			err = errors.Errorf("\tUnknown packet ID %v, %v", group, groupTypeName)
		}

		if err != nil {
			return nil, err
		}
	}

	return &currentStatus, nil
}

func decodeTimestamp(groupTypeName xdi, packet packet, currentStatus *Data) error {
	switch groupTypeName {
	case 0x10: // Utc time
		err := binary.Read(bytes.NewReader(packet.data), binary.BigEndian, &currentStatus.UTCTimestamp)
		return errors.Wrap(err, "could not read utc time")
	case 0x20: // packet counter
		err := binary.Read(bytes.NewReader(packet.data), binary.BigEndian, &currentStatus.PacketCounter)
		return errors.Wrap(err, "could not read Packet counter")
	}
	return errors.Errorf("Unhandled group type name='%v'", groupTypeName)
}

func decodePosition(groupTypeName xdi, packet packet, currentStatus *Data) error {
	switch groupTypeName {
	case 0x20: // Altitude Ellipsoid
		err := binary.Read(bytes.NewReader(packet.data), binary.BigEndian, &currentStatus.AltitudeMEllips)
		return errors.Wrap(err, "could not read ellipsoid altitude")
	case 0x40: // Latitude & Longitude
		err := binary.Read(bytes.NewReader(packet.data), binary.BigEndian, &currentStatus.Latlng)
		return errors.Wrap(err, "could not read (lat,lon)")
	}
	return errors.Errorf("Unhandled group type name='%v'", groupTypeName)
}

func decodeOrientation(groupTypeName xdi, packet packet, currentStatus *Data) error {
	switch groupTypeName {
	case 0x30: // Euler Angles
		err := binary.Read(bytes.NewReader(packet.data), binary.BigEndian, &currentStatus.Euler)
		return errors.Wrap(err, "could not read Euler angles")
	}
	return errors.Errorf("Unhandled group type name='%v'", groupTypeName)
}

func decodeVelocity(groupTypeName xdi, packet packet, currentStatus *Data) error {
	switch groupTypeName {
	case 0x10: // Velocity XYZ
		err := binary.Read(bytes.NewReader(packet.data), binary.BigEndian, &currentStatus.Vel)
		return errors.Wrap(err, "could not read velocity")
	}
	return errors.Errorf("Unhandled group type name='%v'", groupTypeName)
}

func decodeStatusWord(groupTypeName xdi, packet packet, currentStatus *Data) error {
	switch groupTypeName {
	case 0x20: // Status Word
		err := binary.Read(bytes.NewReader(packet.data), binary.BigEndian, &currentStatus.StatusWord)
		return errors.Wrap(err, "could not read status word")
	}
	return errors.Errorf("Unhandled group type name='%v'", groupTypeName)
}

func decodeAcceleration(groupTypeName xdi, packet packet, currentStatus *Data) error {
	switch groupTypeName {
	case 0x10: // DeltaV
		err := binary.Read(bytes.NewReader(packet.data), binary.BigEndian, &currentStatus.DeltaV)
		return errors.Wrap(err, "could not read DeltaV")
	case 0x20: // Acceleration
		err := binary.Read(bytes.NewReader(packet.data), binary.BigEndian, &currentStatus.Acc)
		return errors.Wrap(err, "could not read acceleration")
	case 0x30: // Free Acceleration
		err := binary.Read(bytes.NewReader(packet.data), binary.BigEndian, &currentStatus.FreeAcc)
		return errors.Wrap(err, "could not read free acceleration")
	}
	return errors.Errorf("Unhandled group type name='%v'", groupTypeName)
}

func decodeMagnetic(groupTypeName xdi, packet packet, currentStatus *Data) error {
	switch groupTypeName {
	case 0x20: // MagneticField
		err := binary.Read(bytes.NewReader(packet.data), binary.BigEndian, &currentStatus.Mag)
		return errors.Wrap(err, "could not read magnetic field")
	}
	return errors.Errorf("Unhandled group type name='%v'", groupTypeName)
}

func decodeAngularVelocity(groupTypeName xdi, packet packet, currentStatus *Data) error {
	switch groupTypeName {
	case 0x20: // Rate of turn
		err := binary.Read(bytes.NewReader(packet.data), binary.BigEndian, &currentStatus.AngularVel)
		return errors.Wrap(err, "could not read rate of turn")
	case 0x30: // DeltaQ
		err := binary.Read(bytes.NewReader(packet.data), binary.BigEndian, &currentStatus.DeltaQ)
		return errors.Wrap(err, "could not read DeltaQ")
	}
	return errors.Errorf("Unhandled group type name='%v'", groupTypeName)
}
