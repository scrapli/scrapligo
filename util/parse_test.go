package util_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	scrapligoutil "github.com/scrapli/scrapligo/util"
)

func testTextFsmParse(testName string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("%s: starting", testName)

		s := string(readFile(t, fmt.Sprintf("%s.txt", testName)))

		actualOut, err := scrapligoutil.TextFsmParse(
			s,
			"test-fixtures/cisco_ios_show_version.textfsm",
		)
		if err != nil {
			t.Fatalf("%s: encountered error parsing with textfsm, error: %s", testName, err)
		}

		actualOutJSON, err := json.Marshal(actualOut)
		if err != nil {
			t.Fatalf(
				"%s: encountered error dumping textfsm output to json, error: %s",
				testName,
				err,
			)
		}

		if *update {
			writeGolden(t, testName, actualOutJSON)
		}

		expectedOut := readFile(t, fmt.Sprintf("golden/%s-out.txt", testName))

		if !bytes.Equal(actualOutJSON, expectedOut) {
			t.Fatalf(
				"%s: actual and expected outputs do not match\nactual: %s\nexpected:%s",
				testName,
				actualOut,
				expectedOut,
			)
		}
	}
}

func TestTextFsmParse(t *testing.T) {
	cases := map[string]struct {
		description string
	}{
		"textfsm-parse-simple": {},
	}

	for testName := range cases {
		f := testTextFsmParse(testName)

		t.Run(testName, f)
	}
}
