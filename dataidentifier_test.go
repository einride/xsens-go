package xsens

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
)

func TestDataIdentifier_GetSetUint16(t *testing.T) {
	for _, tt := range []struct {
		dataIdentifier DataIdentifier
		value          uint16
	}{
		{
			dataIdentifier: DataIdentifier{
				DataType:         DataTypeQuaternion,
				Precision:        PrecisionFloat64,
				CoordinateSystem: CoordinateSystemNorthWestUp,
			},
			value: 0x201b,
		},
	} {
		t.Run(fmt.Sprintf("%#x", tt.value), func(t *testing.T) {
			t.Run("Uint16", func(t *testing.T) {
				assert.Equal(t, tt.value, tt.dataIdentifier.Uint16())
			})
			t.Run("SetUint16", func(t *testing.T) {
				var actual DataIdentifier
				actual.SetUint16(tt.value)
				assert.DeepEqual(t, tt.dataIdentifier, actual)
			})
		})
	}
}

func TestDataIdentifier_DataSize(t *testing.T) {
	for _, tt := range []struct {
		dataIdentifier DataIdentifier
		dataSize       uint8
	}{
		{
			dataIdentifier: DataIdentifier{
				DataType:  DataTypeTemperature,
				Precision: PrecisionFloat64,
			},
			dataSize: 8,
		},
		{
			dataIdentifier: DataIdentifier{
				DataType:  DataTypeAcceleration,
				Precision: PrecisionFP1220,
			},
			dataSize: 12,
		},
		{
			dataIdentifier: DataIdentifier{
				DataType:  DataTypeQuaternion,
				Precision: PrecisionFloat32,
			},
			dataSize: 16,
		},
	} {
		t.Run(tt.dataIdentifier.String(), func(t *testing.T) {
			assert.Equal(t, tt.dataSize, tt.dataIdentifier.DataSize())
		})
	}
}
