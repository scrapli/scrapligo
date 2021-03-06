package netconf_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/scrapli/scrapligo/util"
)

func testCommitDiscard(testName string, testCase *util.PayloadTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		d, fileTransportObj := prepareDriver(t, testName, testCase.PayloadFile)

		commitr, err := d.Commit()
		if err != nil {
			t.Fatalf(
				"%s: encountered error running netconf Driver Commit, error: %s",
				testName,
				err,
			)
		}

		if commitr.Failed != nil {
			t.Fatalf("%s: response object indicates failure",
				testName)
		}

		discardr, err := d.Discard()
		if err != nil {
			t.Fatalf(
				"%s: encountered error running netconf Driver Discard, error: %s",
				testName,
				err,
			)
		}

		if discardr.Failed != nil {
			t.Fatalf("%s: response object indicates failure",
				testName)
		}

		actualOut := commitr.Result + "\n" + discardr.Result
		actualIn := bytes.Join(fileTransportObj.Writes, []byte("\n"))

		if *update {
			writeGolden(t, testName, actualIn, actualOut)
		}

		expectedIn := readFile(t, fmt.Sprintf("golden/%s-in.txt", testName))
		expectedOut := readFile(t, fmt.Sprintf("golden/%s-out.txt", testName))

		if !cmp.Equal(actualIn, expectedIn) {
			t.Fatalf(
				"%s: actual and expected inputs do not match\nactual: %s\nexpected:%s",
				testName,
				actualIn,
				expectedIn,
			)
		}

		if !cmp.Equal(actualOut, string(expectedOut)) {
			t.Fatalf(
				"%s: actual and expected outputs do not match\nactual: %s\nexpected:%s",
				testName,
				actualOut,
				expectedOut,
			)
		}
	}
}

func TestCommitDiscard(t *testing.T) {
	cases := map[string]*util.PayloadTestCase{
		"commit-discard-simple": {
			Description: "simple commit/discard test",
			PayloadFile: "commit-discard-simple.txt",
		},
	}

	for testName, testCase := range cases {
		f := testCommitDiscard(testName, testCase)
		t.Run(testName, f)
	}
}
