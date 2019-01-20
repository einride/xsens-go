package xsens

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMessageIdentifier_Ack(t *testing.T) {
	for _, tt := range []struct {
		id  MessageIdentifier
		ack MessageIdentifier
	}{
		{id: MessageIdentifierGotoConfig, ack: MessageIdentifierGotoConfigAck},
		{id: MessageIdentifier(0xfe), ack: MessageIdentifier(0xff)},
	} {
		tt := tt
		t.Run(tt.id.String(), func(t *testing.T) {
			require.Equal(t, tt.ack, tt.id.Ack())
		})
	}
}

func TestMessageIdentifier_IsAck(t *testing.T) {
	require.False(t, MessageIdentifierGotoConfig.IsAck())
	require.True(t, MessageIdentifierGotoConfigAck.IsAck())
}
