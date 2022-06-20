package xsens_test

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	"gotest.tools/v3/assert"
)

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

func requireGoldenFileContent(t *testing.T, goldenFile, actual string) {
	t.Helper()
	goldenFileContent, err := ioutil.ReadFile(goldenFile)
	assert.NilError(t, err)
	expected := string(goldenFileContent)
	if expected != actual {
		diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
			A:        difflib.SplitLines(expected),
			FromFile: "Expected",
			B:        difflib.SplitLines(actual),
			ToFile:   "Actual",
		})
		assert.NilError(t, err)
		t.Fatalf("\nGolden file mismatch:\n%s", diff)
	}
}
