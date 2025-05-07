package xsens

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestMTData2_PacketAt(t *testing.T) {
	for _, tt := range []struct {
		data   MTData2
		i      int
		packet MTData2Packet
	}{
		{
			data:   MTData2{0x12, 0x34, 0x00},
			i:      0,
			packet: MTData2Packet{0x12, 0x34, 0x00},
		},
		{
			data:   MTData2{0x12, 0x34, 0x00, 0x56, 0x78, 0x00},
			i:      3,
			packet: MTData2Packet{0x56, 0x78, 0x00},
		},
	} {
		t.Run(fmt.Sprintf("%v,%v", tt.data, tt.i), func(t *testing.T) {
			packet, err := tt.data.PacketAt(tt.i)
			assert.NilError(t, err)
			assert.DeepEqual(t, tt.packet, packet)
		})
	}
}

func TestMTData2_PacketAt_Error(t *testing.T) {
	for _, tt := range []struct {
		data MTData2
		i    int
	}{
		{data: MTData2{}, i: 0},
		{data: MTData2{0x00}, i: 0},
		{data: MTData2{0x00, 0x00}, i: 0},
		{data: MTData2{0x00, 0x00, 0x01}, i: 0},
	} {
		t.Run(fmt.Sprintf("%v[%v]", tt.data, tt.i), func(t *testing.T) {
			_, err := tt.data.PacketAt(tt.i)
			assert.Assert(t, is.ErrorContains(err, ""))
		})
	}
}

func TestMTData2Packet_Identifier(t *testing.T) {
	assert.DeepEqual(t, DataIdentifier{
		DataType:         DataTypeQuaternion,
		CoordinateSystem: CoordinateSystemNorthEastDown,
		Precision:        PrecisionFP1632,
	}, MTData2Packet{0x20, 0x16, 0x00}.Identifier())
}

func TestMTData2Packet_Data(t *testing.T) {
	for _, tt := range []struct {
		packet MTData2Packet
		data   []byte
	}{
		{
			packet: MTData2Packet{0x12, 0x34, 0x00},
			data:   []byte{},
		},
		{
			packet: MTData2Packet{0x12, 0x34, 0x02, 0x12, 0x34},
			data:   []byte{0x12, 0x34},
		},
	} {
		t.Run(tt.packet.String(), func(t *testing.T) {
			assert.DeepEqual(t, tt.data, tt.packet.Data())
		})
	}
}

func TestIdentifier(t *testing.T) {
	id := DataIdentifier{
		DataType:  DataTypeLatLon,
		Precision: PrecisionFloat32,
	}

	packet := make(MTData2Packet, packetDataStart+id.Precision.Size())
	packet.SetIdentifier(id)
	assert.Equal(t, id, packet.Identifier())
}
