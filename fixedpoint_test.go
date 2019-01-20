package xsens

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFP1220_Float64(t *testing.T) {
	for _, tt := range []struct {
		input    FP1220
		expected float64
	}{
		{
			input:    FP1220{0x00, 0x00, 0x00, 0x00},
			expected: 0,
		},
		{
			input:    FP1220{0x3, 0x9a, 0xdf, 0x3},
			expected: 57.67944622039795,
		},
		{
			input:    FP1220{0x0, 0xbe, 0x1f, 0x81},
			expected: 11.882691383361816,
		},
	} {
		tt := tt
		t.Run(tt.input.String(), func(t *testing.T) {
			require.Equal(t, tt.expected, tt.input.Float64())
		})
	}
}

func TestFP1632_Float64(t *testing.T) {
	for _, tt := range []struct {
		input    FP1632
		expected float64
	}{
		{
			input:    FP1632{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expected: 0,
		},
		{
			input:    FP1632{0xad, 0xf0, 0x52, 0x98, 0x0, 0x39},
			expected: 57.67944828234613,
		},
		{
			input:    FP1632{0xe1, 0xf5, 0x89, 0xe9, 0x0, 0xb},
			expected: 11.882652873406187,
		},
		{
			input:    FP1632{0xf6, 0x3f, 0xca, 0xf0, 0xff, 0xff},
			expected: -0.038089100271463394,
		},
	} {
		tt := tt
		t.Run(tt.input.String(), func(t *testing.T) {
			require.Equal(t, tt.expected, tt.input.Float64())
		})
	}
}
