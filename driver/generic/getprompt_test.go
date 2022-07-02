package generic_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/scrapli/scrapligo/util"

	"github.com/google/go-cmp/cmp"
)

func testGetPrompt(testName string, testCase *util.PayloadTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		d, fileTransportObj := prepareDriver(t, testName, testCase.PayloadFile)

		actualOut, err := d.GetPrompt()
		if err != nil {
			t.Errorf(
				"%s: encountered error running generic Driver GetPrompt, error: %s",
				testName,
				err,
			)
		}

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

func TestGetPrompt(t *testing.T) {
	cases := map[string]*util.PayloadTestCase{
		"get-prompt-simple": {
			Description: "simple get prompt test",
			PayloadFile: "get-prompt-simple.txt",
		},
	}

	for testName, testCase := range cases {
		f := testGetPrompt(testName, testCase)

		t.Run(testName, f)
	}
}
