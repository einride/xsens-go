package xsens_test

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/einride/xsens-go"
	"github.com/stretchr/testify/require"
)

func TestNewMessage_TestData(t *testing.T) {
	for _, tt := range []struct {
		inputFile string
	}{
		{inputFile: "testdata/1/output.bin"},
		{inputFile: "testdata/2/output.bin"},
		{inputFile: "testdata/3/output.bin"},
		{inputFile: "testdata/4/output.bin"},
		{inputFile: "testdata/5/output.bin"},
	} {
		tt := tt
		t.Run(tt.inputFile, func(t *testing.T) {
			f, err := os.Open(tt.inputFile)
			require.NoError(t, err)
			defer func() {
				require.NoError(t, f.Close())
			}()
			sc := bufio.NewScanner(f)
			sc.Split(xsens.ScanMessages)
			for sc.Scan() {
				msg := xsens.Message(sc.Bytes())
				newMsg := xsens.NewMessage(msg.Identifier(), msg.Data())
				require.NoError(t, newMsg.Validate())
				require.Equal(t, msg, newMsg)
			}
			require.NoError(t, sc.Err())
		})
	}
}

func TestNewMessage(t *testing.T) {
	for _, tt := range []struct {
		actual   xsens.Message
		expected xsens.Message
	}{
		{
			actual:   xsens.NewMessage(xsens.MessageIdentifierGotoConfig, nil),
			expected: xsens.Message{0xfa, 0xff, 0x30, 0x00, 0xd1},
		},
		{
			actual:   xsens.NewMessage(xsens.MessageIdentifierError, []byte{byte(xsens.ErrorCodeInvalidMessage)}),
			expected: xsens.Message{0xfa, 0xff, 0x42, 0x01, 0x04, 0xba},
		},
	} {
		tt := tt
		t.Run(fmt.Sprintf("%v", tt.actual), func(t *testing.T) {
			require.Equal(t, tt.expected, tt.actual)
			t.Run("Validate", func(t *testing.T) {
				require.NoError(t, tt.actual.Validate())
			})
		})
	}
}

func TestMessage_Validate_Error(t *testing.T) {
	for _, tt := range []xsens.Message{
		{},
		{0x00},
		{0x00, 0x01},
	} {
		tt := tt
		t.Run(tt.String(), func(t *testing.T) {
			require.Error(t, tt.Validate())
		})
	}
}
