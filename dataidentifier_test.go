package xsens

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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
		tt := tt
		t.Run(fmt.Sprintf("%#x", tt.value), func(t *testing.T) {
			t.Run("Uint16", func(t *testing.T) {
				require.Equal(t, tt.value, tt.dataIdentifier.Uint16())
			})
			t.Run("SetUint16", func(t *testing.T) {
				var actual DataIdentifier
				actual.SetUint16(tt.value)
				require.Equal(t, tt.dataIdentifier, actual)
			})
		})
	}
}
