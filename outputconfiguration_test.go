package xsens_test

import (
	"os"
	"testing"

	"go.einride.tech/xsens"
	"gotest.tools/v3/assert"
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
		t.Run(tt.inputFile, func(t *testing.T) {
			golden, err := os.ReadFile(tt.inputFile)
			assert.NilError(t, err)
			var outputConfiguration xsens.OutputConfiguration
			assert.NilError(t, outputConfiguration.Unmarshal(golden))
			actual, err := outputConfiguration.Marshal()
			assert.NilError(t, err)
			assert.DeepEqual(t, golden, actual)
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
		t.Run(tt.inputFile, func(t *testing.T) {
			input, err := os.ReadFile(tt.inputFile)
			assert.NilError(t, err)
			var outputConfiguration xsens.OutputConfiguration
			assert.NilError(t, outputConfiguration.Unmarshal(input))
			txt, err := outputConfiguration.MarshalText()
			assert.NilError(t, err)
			if shouldUpdateGoldenFiles() {
				assert.NilError(t, os.WriteFile(tt.goldenFile, []byte(txt), 0o600))
			}
			requireGoldenFileContent(t, tt.goldenFile, txt)
		})
	}
}
