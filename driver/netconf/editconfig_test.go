package netconf_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/scrapli/scrapligo/util"
)

func testEditConfig(testName string, testCase *util.PayloadTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		d, fileTransportObj := prepareDriver(t, testName, testCase.PayloadFile)

		r, err := d.EditConfig(
			"candidate",
			"<config>\n    <cdp xmlns=\"http://cisco.com/ns/yang/Cisco-IOS-XR-cdp-cfg\">\n        <timer>80</timer>\n        <enable>true</enable>\n        <log-adjacency></log-adjacency>\n        <hold-time>200</hold-time>\n        <advertise-v1-only></advertise-v1-only>\n    </cdp>\n</config>", //nolint: lll
		)
		if err != nil {
			t.Fatalf(
				"%s: encountered error running network Driver Get, error: %s",
				testName,
				err,
			)
		}

		if r.Failed != nil {
			t.Fatalf("%s: response object indicates failure",
				testName)
		}

		actualOut := r.Result
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

func TestEditConfig(t *testing.T) {
	cases := map[string]*util.PayloadTestCase{
		"edit-config-simple": {
			Description: "simple edit config test",
			PayloadFile: "edit-config-simple.txt",
		},
	}

	for testName, testCase := range cases {
		f := testEditConfig(testName, testCase)
		t.Run(testName, f)
	}
}
