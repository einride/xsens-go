package xsens_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/einride/xsens-go"
	"github.com/stretchr/testify/require"
)

func TestMeasurementData_UnmarshalMessage_TestData(t *testing.T) {
	for _, tt := range []struct {
		inputFile  string
		goldenFile string
	}{
		{inputFile: "testdata/1/output.bin", goldenFile: "testdata/1/output.golden"},
		{inputFile: "testdata/2/output.bin", goldenFile: "testdata/2/output.golden"},
		{inputFile: "testdata/3/output.bin", goldenFile: "testdata/3/output.golden"},
		{inputFile: "testdata/4/output.bin", goldenFile: "testdata/4/output.golden"},
		{inputFile: "testdata/5/output.bin", goldenFile: "testdata/5/output.golden"},
	} {
		tt := tt
		t.Run(tt.inputFile, func(t *testing.T) {
			f, err := os.Open(tt.inputFile)
			require.NoError(t, err)
			defer func() {
				require.NoError(t, f.Close())
			}()
			s := xsens.NewMessageScanner(f)
			var measurementData xsens.MeasurementData
			var actual bytes.Buffer
			for s.Scan() {
				msg := s.Message()
				require.NoError(t, msg.Validate())
				require.NoError(t, measurementData.UnmarshalMTData2(msg))
				text, err := measurementData.MarshalText()
				require.NoError(t, err)
				_, err = actual.WriteString(text)
				require.NoError(t, err)
				require.NoError(t, actual.WriteByte('\n'))
				require.NoError(t, actual.WriteByte('\n'))
			}
			if shouldUpdateGoldenFiles() {
				require.NoError(t, ioutil.WriteFile(tt.goldenFile, actual.Bytes(), 0644))
			}
			requireGoldenFileContent(t, tt.goldenFile, actual.String())
		})
	}
}
