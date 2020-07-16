package xsens

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestDataType_HasOutputFormat(t *testing.T) {
	assert.Assert(t, DataTypeEulerAngles.HasPrecision())
	assert.Assert(t, !DataTypePacketCounter.HasPrecision())
}
