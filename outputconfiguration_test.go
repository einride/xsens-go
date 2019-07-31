package xsens_test

import (
	"io/ioutil"
	"testing"

	"github.com/einride/xsens-go"
	"github.com/stretchr/testify/require"
)

func TestOutputConfiguration_UnmarshalMarshal_TestData(t *testing.T) {
	for _, tt := range []struct {
		inputFile string
	}{
		{inputFile: "testdata/1/outputconfig.bin"},
		{inputFile: "testdata/2/outputconfig.bin"},
		{inputFile: "testdata/3/outputconfig.bin"},
		{inputFile: "testdata/4/outputconfig.bin"},
		{inputFile: "testdata/5/outputconfig.bin"},
	} {
		tt := tt
		t.Run(tt.inputFile, func(t *testing.T) {
			golden, err := ioutil.ReadFile(tt.inputFile)
			require.NoError(t, err)
			var outputConfiguration xsens.OutputConfiguration
			require.NoError(t, outputConfiguration.Unmarshal(golden))
			actual, err := outputConfiguration.Marshal()
			require.NoError(t, err)
			require.Equal(t, golden, actual)
		})
	}
}

func TestOutputConfiguration_MarshalText_TestData(t *testing.T) {
	for _, tt := range []struct {
		inputFile  string
		goldenFile string
	}{
		{inputFile: "testdata/1/outputconfig.bin", goldenFile: "testdata/1/outputconfig.golden"},
		{inputFile: "testdata/2/outputconfig.bin", goldenFile: "testdata/2/outputconfig.golden"},
		{inputFile: "testdata/3/outputconfig.bin", goldenFile: "testdata/3/outputconfig.golden"},
		{inputFile: "testdata/4/outputconfig.bin", goldenFile: "testdata/4/outputconfig.golden"},
		{inputFile: "testdata/5/outputconfig.bin", goldenFile: "testdata/5/outputconfig.golden"},
	} {
		tt := tt
		t.Run(tt.inputFile, func(t *testing.T) {
			input, err := ioutil.ReadFile(tt.inputFile)
			require.NoError(t, err)
			var outputConfiguration xsens.OutputConfiguration
			require.NoError(t, outputConfiguration.Unmarshal(input))
			txt, err := outputConfiguration.MarshalText()
			require.NoError(t, err)
			if shouldUpdateGoldenFiles() {
				require.NoError(t, ioutil.WriteFile(tt.goldenFile, []byte(txt), 0644))
			}
			requireGoldenFileContent(t, tt.goldenFile, txt)
		})
	}
}