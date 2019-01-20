package xsens_test

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	_ = flag.Bool(updateGoldenFilesFlag())
	os.Exit(m.Run())
}

func updateGoldenFilesFlag() (string, bool, string) {
	return "update", false, "Update golden files."
}

func shouldUpdateGoldenFiles() bool {
	var flags flag.FlagSet
	flags.SetOutput(ioutil.Discard)
	update := flags.Bool(updateGoldenFilesFlag())
	_ = flags.Parse(os.Args[1:]) // error will always be an unparsed flags error
	return *update
}

func requireGoldenFileContent(t *testing.T, goldenFile string, actual string) {
	t.Helper()
	goldenFileContent, err := ioutil.ReadFile(goldenFile)
	require.NoError(t, err)
	expected := string(goldenFileContent)
	if expected != actual {
		diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
			A:        difflib.SplitLines(expected),
			FromFile: "Expected",
			B:        difflib.SplitLines(actual),
			ToFile:   "Actual",
		})
		require.NoError(t, err)
		t.Fatalf("\nGolden file mismatch:\n%s", diff)
	}
}
