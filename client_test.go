package xsens

import (
	"bytes"
	"io"
	"net"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var inputDataWithHeader = []byte{
	// Header
	0xFA, 0xFF, 54,
	// Indication that this message is extended (>254)
	0xFF,
	// Length of message
	0x1, 0x19,
	// Timestamp
	0x10, 0x10, 0x0C, 0x29, 0xFA, 0xC3, 0xE0, 0x07, 0xE2, 0x08, 0x09, 0x0D, 0x3B, 0x04, 0xF7,
	// Packetcounter
	0x10, 0x20, 0x02, 0x43, 0x80,
	// Euler angles
	0x20, 0x33, 0x18, 0x3F, 0xB6, 0x81, 0x2E, 0x60, 0x00, 0x00, 0x00, 0x40, 0x15, 0x3A, 0x60, 0x80, 0x00, 0x00, 0x00,
	0xC0, 0x06, 0x41, 0x85, 0xA0, 0x00, 0x00, 0x00,
	// DeltaV
	0x40, 0x13, 0x18, 0xBF, 0x82, 0x88, 0x0F, 0x20, 0x00, 0x00, 0x00, 0x3F, 0x28, 0xAD, 0xFA, 0x00, 0x00, 0x00, 0x00,
	0x3F, 0xB9, 0x0D, 0xF2, 0x40, 0x00, 0x00, 0x00,
	// Acceleration
	0x40, 0x23, 0x18, 0xBF, 0xEC, 0xF4, 0x42, 0x80, 0x00, 0x00, 0x00, 0x3F, 0x93, 0x3E, 0x01, 0xC0, 0x00, 0x00, 0x00,
	0x40, 0x23, 0x92, 0xE5, 0xE0, 0x00, 0x00, 0x00,
	// Free Acceleration
	0x40, 0x33, 0x18, 0x3F, 0x48, 0x9E, 0x3F, 0x00, 0x00, 0x00, 0x00, 0x3F, 0x71, 0xFA, 0x9D, 0xC0, 0x00, 0x00, 0x00,
	0x3F, 0x8B, 0xCA, 0xA0, 0x00, 0x00, 0x00, 0x00,
	// Altitude Ellipsoid
	0x50, 0x23, 0x08, 0x40, 0x59, 0x8F, 0x1D, 0x5D, 0x2D, 0x99, 0xF1,
	// Latitude Longitude
	0x50, 0x43, 0x10, 0x40, 0x45, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x2a, 0xbd, 0x70, 0xa3, 0xd7, 0x0a, 0x3d,
	// Rate of turn
	0x80, 0x23, 0x18, 0xBF, 0x40, 0xA0, 0x4B, 0xE0, 0x00, 0x00, 0x01, 0xBF, 0x3B, 0x06, 0xBF, 0x00, 0x00, 0x00, 0x00,
	0x3F, 0x55, 0x72, 0x9E, 0x80, 0x00, 0x00, 0x00,
	// DeltaQ
	0x80, 0x33, 0x20, 0x3F, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xBE, 0xC5, 0x48, 0x0F, 0x00, 0x00, 0x00, 0x00,
	0xBE, 0xC1, 0x4B, 0xFF, 0x40, 0x00, 0x00, 0x00, 0x3E, 0xDB, 0x73, 0xFD, 0xE0, 0x00, 0x00, 0x00,
	// Magnetometer
	0xC0, 0x23, 0x18, 0x3F, 0xC8, 0x6A, 0x38, 0xE0, 0x00, 0x00, 0x00, 0x3F, 0xD4, 0x09, 0xF6, 0x80, 0x00, 0x00, 0x00,
	0xBF, 0xF3, 0x4D, 0x1D, 0xA0, 0x00, 0x00, 0x00,
	// Velocity
	0xD0, 0x13, 0x18, 0x3F, 0xB2, 0x99, 0x03, 0xC0, 0x00, 0x00, 0x00, 0x3F, 0x87, 0x31, 0xDC, 0x40, 0x00, 0x00, 0x00,
	0x3F, 0x92, 0x86, 0x15, 0x80, 0x00, 0x00, 0x00,
	// Statusword
	0xE0, 0x20, 0x04, 0x01, 0x80, 0x00, 0x47,
	// Checksum
	0x00,
}

// A GNSS contains UTCtime and Packetcounter as well as the GNSS data
var gnssDataWithHeader = []byte{
	// Header
	0xFA, 0xFF, 54,
	0x75,
	// Timestamp
	0x10, 0x10, 0x0C, 0x1D, 0xCD, 0x65, 0x00, 0x07, 0xE2, 0x08, 0x09, 0x0D, 0x3B, 0x02, 0xF7,
	// Packetcounter
	0x10, 0x20, 0x02, 0x42, 0x97,
	// GNSS message
	0x70, 0x20, 0x78, 0x17, 0x99, 0xE0, 0xB4, 0x1C, 0x00, 0x00, 0x00, 0x00, 0x01, 0x15, 0x04, 0x00, 0x08, 0x1F, 0x17,
	0x00, 0x0A, 0x30, 0x1F, 0x00, 0x0B, 0x1C, 0x1F, 0x00, 0x0E, 0x1D, 0x06, 0x00, 0x12, 0x20, 0x1F, 0x00, 0x14, 0x2C,
	0x1F, 0x00, 0x16, 0x1C, 0x1F, 0x00, 0x18, 0x21, 0x1F, 0x00, 0x1B, 0x1C, 0x1F, 0x00, 0x1C, 0x1C, 0x1F, 0x00, 0x20,
	0x2D, 0x1F, 0x01, 0x78, 0x00, 0x01, 0x01, 0x7C, 0x00, 0x01, 0x01, 0x7E, 0x00, 0x01, 0x05, 0x01, 0x00, 0x01, 0x05,
	0x02, 0x00, 0x01, 0x05, 0x05, 0x00, 0x01, 0x06, 0x04, 0x23, 0x1F, 0x06, 0x05, 0x00, 0x21, 0x06, 0x06, 0x15, 0x14,
	0x06, 0x0B, 0x17, 0x14, 0x06, 0x0C, 0x1B, 0x16, 0x06, 0x0D, 0x22, 0x17, 0x06, 0x14, 0x29, 0x1F, 0x06, 0x15, 0x1E,
	0x1E, 0x06, 0x16, 0x1C, 0x17, 0x06, 0xFF, 0x2E, 0x07}

func TestReadGNSS(t *testing.T) {
	var data Data
	assert.Equal(t, io.EOF, errors.Cause(data.Read(bytes.NewReader(gnssDataWithHeader))))
}

func TestSkipGNSS(t *testing.T) {
	// Data should only be set from non-GNSS message
	reader, writer := net.Pipe()
	go func() {
		// GNSS message, should not set data and will get EOF error
		_, err := writer.Write(gnssDataWithHeader)
		// Regular message, should set data
		assert.Nil(t, err)
		_, err = writer.Write(inputDataWithHeader)
		assert.Nil(t, err)
		// GNSS message, should not set data and will get EOF error
		_, err = writer.Write(gnssDataWithHeader)
		assert.Nil(t, err)
		err = writer.Close()
		assert.Nil(t, err)
	}()
	var data Data
	assert.NoError(t, data.Read(reader))
	assert.Equal(t, 13.37, data.Latlng.Lng)
	assert.Equal(t, 42.0, data.Latlng.Lat)
	assert.Equal(t, io.EOF, errors.Cause(data.Read(reader)))
}

func TestReadMsgs(t *testing.T) {
	reader, writer := net.Pipe()
	go func() {
		_, err := writer.Write(inputDataWithHeader)
		assert.Nil(t, err)
		err = writer.Close()
		assert.Nil(t, err)
	}()
	var data Data
	assert.NoError(t, data.Read(reader))
	assert.Equal(t, 42.0, data.Latlng.Lat)
	assert.Equal(t, 13.37, data.Latlng.Lng)
	assert.Equal(t, io.EOF, errors.Cause(data.Read(reader)))
}

func TestReadTwoMsgs(t *testing.T) {
	reader, writer := net.Pipe()
	go func() {
		_, err := writer.Write(inputDataWithHeader)
		assert.Nil(t, err)
		_, err = writer.Write(inputDataWithHeader)
		assert.Nil(t, err)
		err = writer.Close()
		assert.Nil(t, err)
	}()
	var data Data
	assert.NoError(t, data.Read(reader))
	assert.Equal(t, 42.0, data.Latlng.Lat)
	assert.Equal(t, 13.37, data.Latlng.Lng)
	assert.NoError(t, data.Read(reader))
	assert.Equal(t, 42.0, data.Latlng.Lat)
	assert.Equal(t, 13.37, data.Latlng.Lng)
	assert.Equal(t, io.EOF, errors.Cause(data.Read(reader)))
}

func TestReadNextHeader(t *testing.T) {
	foundHeaders := []struct {
		name  string
		bytes []byte
	}{
		{
			name: "Normal case",
			bytes: []byte{
				// Header start here
				0xFA,
				0xFF,
				54,
				0x2A, // length
			},
		},
		{
			name: "Noise before header",
			bytes: []byte{
				0x00,
				0x01,
				0x02,
				// Header start here
				0xFA,
				0xFF,
				54,
				0x2A, // length
			},
		},
		{
			name: "Unhandled packet before header",
			bytes: []byte{
				// Header of unhandled packet starts here
				0xFA,
				0x42,
				54,
				0x04,
				// Data of unhandled packet (does not match length?! we don't care.)
				0x03,
				0x02,
				0x01,
				// Header start here
				0xFA,
				0xFF,
				54,
				0x2A, // length
			},
		},
	}
	for _, tc := range foundHeaders {
		t.Run(tc.name, func(t *testing.T) {
			reader, writer := net.Pipe()
			go func() {
				_, err := writer.Write(tc.bytes)
				assert.Nil(t, err)
				err = writer.Close()
				assert.Nil(t, err)
			}()
			var h header
			assert.NoError(t, h.Read(reader))
			assert.Equal(t, uint8(42), h.LEN)
		})
	}
}

func TestReadNextHeader_Error_EOF(t *testing.T) {
	foundHeaders := []struct {
		bytes    []byte
		expected error
	}{
		{bytes: []byte{}, expected: io.EOF},
		{bytes: []byte{0xFA}, expected: io.EOF},
		// Because we get EOF in the middle of a read it's called an unexpected EOF, as opposed to before reading as in
		// the prior cases
		{bytes: []byte{0xFA, 0xFF}, expected: io.ErrUnexpectedEOF},
		{bytes: []byte{0xFA, 0xFF, 54}, expected: io.ErrUnexpectedEOF},
	}
	for _, tc := range foundHeaders {
		t.Run("", func(t *testing.T) {
			reader, writer := net.Pipe()
			go func() {
				_, err := writer.Write(tc.bytes)
				assert.Nil(t, err)
				err = writer.Close()
				assert.Nil(t, err)
			}()
			var h header
			assert.Equal(t, tc.expected, errors.Cause(h.Read(reader)))
		})
	}
}
