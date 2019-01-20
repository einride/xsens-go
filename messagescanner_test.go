package xsens_test

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"os"
	"testing"

	"github.com/einride/xsens-go"
	"github.com/stretchr/testify/require"
)

func TestMessageScanner_Scan_TestData(t *testing.T) {
	for _, tt := range []struct {
		inputFile  string
		goldenFile string
	}{
		{inputFile: "testdata/1/output.bin", goldenFile: "testdata/1/messages.golden"},
		{inputFile: "testdata/2/output.bin", goldenFile: "testdata/2/messages.golden"},
		{inputFile: "testdata/3/output.bin", goldenFile: "testdata/3/messages.golden"},
		{inputFile: "testdata/4/output.bin", goldenFile: "testdata/4/messages.golden"},
		{inputFile: "testdata/5/output.bin", goldenFile: "testdata/5/messages.golden"},
	} {
		tt := tt
		t.Run(tt.inputFile, func(t *testing.T) {
			f, err := os.Open(tt.inputFile)
			require.NoError(t, err)
			defer func() {
				require.NoError(t, f.Close())
			}()
			var actual bytes.Buffer
			s := xsens.NewMessageScanner(f)
			for s.Scan() {
				m := s.Message()
				require.NoError(t, m.Validate())
				_, err = actual.WriteString(hex.EncodeToString(m))
				require.NoError(t, err)
				require.NoError(t, actual.WriteByte('\n'))
			}
			require.NoError(t, s.Err())
			if shouldUpdateGoldenFiles() {
				require.NoError(t, ioutil.WriteFile(tt.goldenFile, actual.Bytes(), 0644))
			}
			requireGoldenFileContent(t, tt.goldenFile, actual.String())
		})
	}
}
