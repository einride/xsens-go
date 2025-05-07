package xsens

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestMessageIdentifier_Ack(t *testing.T) {
	for _, tt := range []struct {
		id  MessageIdentifier
		ack MessageIdentifier
	}{
		{id: MessageIdentifierGotoConfig, ack: MessageIdentifierGotoConfigAck},
		{id: MessageIdentifier(0xfe), ack: MessageIdentifier(0xff)},
	} {
		t.Run(tt.id.String(), func(t *testing.T) {
			assert.Equal(t, tt.ack, tt.id.Ack())
		})
	}
}

func TestMessageIdentifier_IsAck(t *testing.T) {
	assert.Assert(t, !MessageIdentifierGotoConfig.IsAck())
	assert.Assert(t, MessageIdentifierGotoConfigAck.IsAck())
}
