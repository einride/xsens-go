package xsens

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDataType_HasOutputFormat(t *testing.T) {
	require.True(t, DataTypeEulerAngles.HasPrecision())
	require.False(t, DataTypePacketCounter.HasPrecision())
}
