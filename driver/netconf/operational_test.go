package netconf_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/scrapli/scrapligo/util"
)

func testGetOperational(testName string, testCase *operationalTest) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		d, fileTransportObj := prepareDriver(t, testName, testCase.payload.PayloadFile)

		r, err := d.Get(testCase.filter)
		if err != nil {
			t.Fatalf(
				"%s: encountered error running netconf Driver Validate, error: %s",
				testName,
				err,
			)
		}

		if r.Failed != nil && !testCase.shouldFail {
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
			t.Errorf(
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

type operationalTest struct {
	payload    *util.PayloadTestCase
	filter     string
	shouldFail bool
}

func TestGetOperational(t *testing.T) {
	cases := map[string]*operationalTest{
		"junos-getoperational-components-fail": {
			payload: &util.PayloadTestCase{
				Description: "operational show components",
				PayloadFile: "junos-getoperational-components-fail.txt",
			},
			filter: `<components xmlns="http://openconfig.net/yang/platform">
			<component>
			<state>
			<type xmlns:idx="http://openconfig.net/yang/platform-types">idx:CHASSIS</type>
			</state>
			</component>
			</components>`,
			shouldFail: true,
		},
	}

	for testName, testCase := range cases {
		f := testGetOperational(testName, testCase)
		t.Run(testName, f)
	}
}
