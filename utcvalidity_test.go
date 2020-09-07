package xsens

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestIsDateValid(t *testing.T) {
	invalid := UTCValidity(0)
	assert.Assert(t, !invalid.IsDateValid())

	valid := UTCValidity(0 | UTCDateValidFlag)
	assert.Assert(t, valid.IsDateValid())
}

func TestIsTimeOfDayValid(t *testing.T) {
	invalid := UTCValidity(0)
	assert.Assert(t, !invalid.IsTimeOfDayValid())

	valid := UTCValidity(0 | UTCTimeOfDayValidFlag)
	assert.Assert(t, valid.IsTimeOfDayValid())
}

func TestIsTimeOfDayFullyResolved(t *testing.T) {
	invalid := UTCValidity(0)
	assert.Assert(t, !invalid.IsTimeOfDayFullyResolved())

	valid := UTCValidity(0 | UTCTimeOfDayFullyResolvedFlag)
	assert.Assert(t, valid.IsTimeOfDayFullyResolved())
}
