package xsens

import (
	"bytes"
	"encoding/binary"
)

// ScanMessages is a bufio.Scanner SplitFunc that splits a stream into Xsens messages.
func ScanMessages(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) == 0 && atEOF {
		return 0, nil, nil
	}
	// scan to next message header index
	indexOfMessage := bytes.Index(data, []byte{valueOfPreamble, valueOfBusIdentifier})
	if indexOfMessage == -1 {
		if data[len(data)-1] == valueOfPreamble {
			// edge case: latest byte is a preamble, so keep it
			return len(data) - 1, nil, nil
		}
		// found no message header, discard all data and read more
		return len(data), nil, nil
	}
	message := data[indexOfMessage:]
	if len(message) < indexOfLength+lengthOfLength {
		// not enough data to read length: advance to message start and read more data
		return indexOfMessage, nil, nil
	}
	indexOfData := indexOfData
	lengthOfData := int(message[indexOfLength])
	if lengthOfData == valueOfLengthExtended {
		if len(message) < indexOfExtendedLength+lengthOfExtendedLength {
			// not enough data to read extended length: advance to message start and read more data
			return indexOfMessage, nil, nil
		}
		indexOfData = indexOfExtendedData
		lengthOfData = int(binary.BigEndian.Uint16(
			message[indexOfExtendedLength : indexOfExtendedLength+lengthOfExtendedLength]))
	}
	lengthOfMessage := indexOfData + lengthOfData + lengthOfChecksum
	if len(message) < lengthOfMessage {
		// not enough data to read data and checksum, advance to message start and read more data
		return indexOfMessage, nil, nil
	}
	// message now contains the exact bytes of the next scanned message
	message = message[:lengthOfMessage]
	return indexOfMessage + lengthOfMessage, message, nil
}
