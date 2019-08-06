package xsens

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrecision_Size(t *testing.T) {
	for _, tt := range []struct {
		precision Precision
		size      uint8
	}{
		{precision: PrecisionFloat32, size: 4},
		{precision: PrecisionFP1220, size: 4},
		{precision: PrecisionFP1632, size: 6},
		{precision: PrecisionFloat64, size: 8},
	} {
		tt := tt
		t.Run(tt.precision.String(), func(t *testing.T) {
			require.Equal(t, tt.size, tt.precision.Size())
		})
	}
}
