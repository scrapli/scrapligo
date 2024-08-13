package response_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/scrapli/scrapligo/response"
)

type testNetconfRecordTestCase struct {
	description string
	version     string
	payloadFile string
}

func testNetconfRecord(testName string, testCase *testNetconfRecordTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		r := response.NewNetconfResponse(nil, nil, "localhost", 830, testCase.version)

		r.Record(readFile(t, testCase.payloadFile))

		// set the timestamp bits to nil so we dont compare those
		r.StartTime = time.Time{}
		r.EndTime = time.Time{}
		r.ElapsedTime = 0

		actual, err := json.Marshal(r)
		if err != nil {
			t.Fatalf("failed marshaling platform, error: %s", err)
		}

		if *update {
			writeGolden(t, testName, actual)
		}

		expected := readFile(t, fmt.Sprintf("golden/%s-out.txt", testName))

		if !cmp.Equal(actual, expected) {
			t.Fatalf(
				"%s: actual and expected inputs do not match\nactual: %s\nexpected:%s",
				testName,
				actual,
				expected,
			)
		}
	}
}

func TestNetconfRecord(t *testing.T) {
	cases := map[string]*testNetconfRecordTestCase{
		"record-response-10": {
			description: "simple test to test recording netconf 1.0 response",
			version:     "1.0",
			payloadFile: "netconf-output-10.txt",
		},
		"record-response-10-errors": {
			description: "simple test to test recording netconf 1.0 response with errors",
			version:     "1.0",
			payloadFile: "netconf-output-10-errors.txt",
		},
		"record-response-11": {
			description: "simple test to test recording netconf 1.1 response",
			version:     "1.1",
			payloadFile: "netconf-output-11.txt",
		},
		"record-response-11b": {
			description: "a tricky banner to parse test to test recording netconf 1.1 response",
			version:     "1.1",
			payloadFile: "netconf-output-11b.txt",
		},
	}

	for testName, testCase := range cases {
		f := testNetconfRecord(testName, testCase)

		t.Run(testName, f)
	}
}
